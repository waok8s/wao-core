//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EndpointTerm) DeepCopyInto(out *EndpointTerm) {
	*out = *in
	if in.BasicAuthSecret != nil {
		in, out := &in.BasicAuthSecret, &out.BasicAuthSecret
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	if in.FetchInterval != nil {
		in, out := &in.FetchInterval, &out.FetchInterval
		*out = new(metav1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EndpointTerm.
func (in *EndpointTerm) DeepCopy() *EndpointTerm {
	if in == nil {
		return nil
	}
	out := new(EndpointTerm)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsCollector) DeepCopyInto(out *MetricsCollector) {
	*out = *in
	in.InletTemp.DeepCopyInto(&out.InletTemp)
	in.DeltaP.DeepCopyInto(&out.DeltaP)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsCollector.
func (in *MetricsCollector) DeepCopy() *MetricsCollector {
	if in == nil {
		return nil
	}
	out := new(MetricsCollector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfig) DeepCopyInto(out *NodeConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfig.
func (in *NodeConfig) DeepCopy() *NodeConfig {
	if in == nil {
		return nil
	}
	out := new(NodeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfigList) DeepCopyInto(out *NodeConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfigList.
func (in *NodeConfigList) DeepCopy() *NodeConfigList {
	if in == nil {
		return nil
	}
	out := new(NodeConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfigSpec) DeepCopyInto(out *NodeConfigSpec) {
	*out = *in
	in.MetricsCollector.DeepCopyInto(&out.MetricsCollector)
	in.Predictor.DeepCopyInto(&out.Predictor)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfigSpec.
func (in *NodeConfigSpec) DeepCopy() *NodeConfigSpec {
	if in == nil {
		return nil
	}
	out := new(NodeConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfigStatus) DeepCopyInto(out *NodeConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfigStatus.
func (in *NodeConfigStatus) DeepCopy() *NodeConfigStatus {
	if in == nil {
		return nil
	}
	out := new(NodeConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfigTemplate) DeepCopyInto(out *NodeConfigTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfigTemplate.
func (in *NodeConfigTemplate) DeepCopy() *NodeConfigTemplate {
	if in == nil {
		return nil
	}
	out := new(NodeConfigTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeConfigTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfigTemplateList) DeepCopyInto(out *NodeConfigTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeConfigTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfigTemplateList.
func (in *NodeConfigTemplateList) DeepCopy() *NodeConfigTemplateList {
	if in == nil {
		return nil
	}
	out := new(NodeConfigTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeConfigTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfigTemplateSpec) DeepCopyInto(out *NodeConfigTemplateSpec) {
	*out = *in
	in.NodeSelector.DeepCopyInto(&out.NodeSelector)
	in.MetricsCollector.DeepCopyInto(&out.MetricsCollector)
	in.Predictor.DeepCopyInto(&out.Predictor)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfigTemplateSpec.
func (in *NodeConfigTemplateSpec) DeepCopy() *NodeConfigTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(NodeConfigTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfigTemplateStatus) DeepCopyInto(out *NodeConfigTemplateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfigTemplateStatus.
func (in *NodeConfigTemplateStatus) DeepCopy() *NodeConfigTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(NodeConfigTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Predictor) DeepCopyInto(out *Predictor) {
	*out = *in
	if in.PowerConsumption != nil {
		in, out := &in.PowerConsumption, &out.PowerConsumption
		*out = new(EndpointTerm)
		(*in).DeepCopyInto(*out)
	}
	if in.PowerConsumptionEndpointProvider != nil {
		in, out := &in.PowerConsumptionEndpointProvider, &out.PowerConsumptionEndpointProvider
		*out = new(EndpointTerm)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Predictor.
func (in *Predictor) DeepCopy() *Predictor {
	if in == nil {
		return nil
	}
	out := new(Predictor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemplateData) DeepCopyInto(out *TemplateData) {
	*out = *in
	out.IPv4 = in.IPv4
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemplateData.
func (in *TemplateData) DeepCopy() *TemplateData {
	if in == nil {
		return nil
	}
	out := new(TemplateData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemplateDataIPv4) DeepCopyInto(out *TemplateDataIPv4) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemplateDataIPv4.
func (in *TemplateDataIPv4) DeepCopy() *TemplateDataIPv4 {
	if in == nil {
		return nil
	}
	out := new(TemplateDataIPv4)
	in.DeepCopyInto(out)
	return out
}
