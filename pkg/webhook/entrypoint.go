package webhook

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"os"

	"github.com/flyteorg/flytepropeller/pkg/controller/config"
	"github.com/flyteorg/flytepropeller/pkg/controller/executors"
	"github.com/flyteorg/flytepropeller/pkg/utils"
	config2 "github.com/flyteorg/flytepropeller/pkg/webhook/config"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/flyteorg/flytestdlib/promutils"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	PodNameEnvVar      = "POD_NAME"
	PodNamespaceEnvVar = "POD_NAMESPACE"
)

func Run(ctx context.Context, propellerCfg *config.Config, cfg *config2.Config, defaultNamespace string) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Println(string(raw))

	kubeClient, kubecfg, err := utils.GetKubeConfig(ctx, propellerCfg)
	if err != nil {
		return err
	}

	// Add the propeller subscope because the MetricsPrefix only has "flyte:" to get uniform collection of metrics.
	propellerScope := promutils.NewScope(cfg.MetricsPrefix).NewSubScope("propeller").NewSubScope(propellerCfg.LimitNamespace)
	webhookScope := propellerScope.NewSubScope("webhook")

	limitNamespace := ""
	if propellerCfg.LimitNamespace != defaultNamespace {
		limitNamespace = propellerCfg.LimitNamespace
	}

	secretsWebhook := NewPodMutator(cfg, webhookScope)

	// Creates a MutationConfig to instruct ApiServer to call this service whenever a Pod is being created.
	err = createMutationConfig(ctx, kubeClient, secretsWebhook, defaultNamespace)
	if err != nil {
		return err
	}

	mgr, err := manager.New(kubecfg, manager.Options{
		Port:          cfg.ListenPort,
		CertDir:       cfg.CertDir,
		Namespace:     limitNamespace,
		SyncPeriod:    &propellerCfg.DownstreamEval.Duration,
		ClientBuilder: executors.NewFallbackClientBuilder(webhookScope),
	})

	if err != nil {
		logger.Fatalf(ctx, "Failed to initialize controller run-time manager. Error: %v", err)
	}

	err = secretsWebhook.Register(ctx, mgr)
	if err != nil {
		logger.Fatalf(ctx, "Failed to register webhook with manager. Error: %v", err)
	}

	logger.Infof(ctx, "Starting controller-runtime manager")
	return mgr.Start(ctx)
}

func createMutationConfig(ctx context.Context, kubeClient *kubernetes.Clientset, webhookObj *PodMutator, defaultNamespace string) error {
	shouldAddOwnerRef := true
	podName, found := os.LookupEnv(PodNameEnvVar)
	if !found {
		shouldAddOwnerRef = false
	}

	podNamespace, found := os.LookupEnv(PodNamespaceEnvVar)
	if !found {
		shouldAddOwnerRef = false
		podNamespace = defaultNamespace
	}

	mutateConfig, err := webhookObj.CreateMutationWebhookConfiguration(podNamespace)
	if err != nil {
		return err
	}

	if shouldAddOwnerRef {
		// Lookup the pod to retrieve its UID
		p, err := kubeClient.CoreV1().Pods(podNamespace).Get(ctx, podName, v1.GetOptions{})
		if err != nil {
			logger.Infof(ctx, "Failed to get Pod [%v/%v]. Error: %v", podNamespace, podName, err)
			return fmt.Errorf("failed to get pod. Error: %w", err)
		}

		mutateConfig.OwnerReferences = p.OwnerReferences
	}

	logger.Infof(ctx, "Creating MutatingWebhookConfiguration [%v/%v]", mutateConfig.GetNamespace(), mutateConfig.GetName())

	_, err = kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(ctx, mutateConfig, v1.CreateOptions{})
	var statusErr *errors.StatusError
	if err != nil && errors2.As(err, &statusErr) && statusErr.Status().Reason == v1.StatusReasonAlreadyExists {
		logger.Infof(ctx, "Failed to create MutatingWebhookConfiguration. Will attempt to update. Error: %v", err)
		obj, getErr := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(ctx, mutateConfig.Name, v1.GetOptions{})
		if getErr != nil {
			logger.Infof(ctx, "Failed to get MutatingWebhookConfiguration. Error: %v", getErr)
			return err
		}

		obj.Webhooks = mutateConfig.Webhooks
		_, err = kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(ctx, obj, v1.UpdateOptions{})
		if err == nil {
			logger.Infof(ctx, "Successfully updated existing mutating webhook config.")
		}

		return err
	}

	return nil
}
