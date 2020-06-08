// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRequest) DeepCopyInto(out *ResourceRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRequest.
func (in *ResourceRequest) DeepCopy() *ResourceRequest {
	if in == nil {
		return nil
	}
	out := new(ResourceRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRequestList) DeepCopyInto(out *ResourceRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRequestList.
func (in *ResourceRequestList) DeepCopy() *ResourceRequestList {
	if in == nil {
		return nil
	}
	out := new(ResourceRequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRequestSpec) DeepCopyInto(out *ResourceRequestSpec) {
	*out = *in
	if in.SubmittedAt != nil {
		in, out := &in.SubmittedAt, &out.SubmittedAt
		*out = (*in).DeepCopy()
	}
	if in.QueuingBudget != nil {
		in, out := &in.QueuingBudget, &out.QueuingBudget
		*out = new(v1.Duration)
		**out = **in
	}
	out.Cpu = in.Cpu.DeepCopy()
	out.Memory = in.Memory.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRequestSpec.
func (in *ResourceRequestSpec) DeepCopy() *ResourceRequestSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceRequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRequestStatus) DeepCopyInto(out *ResourceRequestStatus) {
	*out = *in
	if in.AcceptedAt != nil {
		in, out := &in.AcceptedAt, &out.AcceptedAt
		*out = (*in).DeepCopy()
	}
	if in.LastUpdatedAt != nil {
		in, out := &in.LastUpdatedAt, &out.LastUpdatedAt
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRequestStatus.
func (in *ResourceRequestStatus) DeepCopy() *ResourceRequestStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceRequestStatus)
	in.DeepCopyInto(out)
	return out
}
