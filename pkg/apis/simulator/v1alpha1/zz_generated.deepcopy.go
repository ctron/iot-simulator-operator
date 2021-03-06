// +build !ignore_autogenerated

/*******************************************************************************
 * Copyright (c) 2018 Red Hat Inc
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0
 *
 * SPDX-License-Identifier: EPL-2.0
 *******************************************************************************/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdapterEndpoints) DeepCopyInto(out *AdapterEndpoints) {
	*out = *in
	out.HTTP = in.HTTP
	out.MQTT = in.MQTT
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdapterEndpoints.
func (in *AdapterEndpoints) DeepCopy() *AdapterEndpoints {
	if in == nil {
		return nil
	}
	out := new(AdapterEndpoints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Build) DeepCopyInto(out *Build) {
	*out = *in
	out.Git = in.Git
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Build.
func (in *Build) DeepCopy() *Build {
	if in == nil {
		return nil
	}
	out := new(Build)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonSpec) DeepCopyInto(out *CommonSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonSpec.
func (in *CommonSpec) DeepCopy() *CommonSpec {
	if in == nil {
		return nil
	}
	out := new(CommonSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsumerSpec) DeepCopyInto(out *ConsumerSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsumerSpec.
func (in *ConsumerSpec) DeepCopy() *ConsumerSpec {
	if in == nil {
		return nil
	}
	out := new(ConsumerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsumerStatus) DeepCopyInto(out *ConsumerStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsumerStatus.
func (in *ConsumerStatus) DeepCopy() *ConsumerStatus {
	if in == nil {
		return nil
	}
	out := new(ConsumerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitSource) DeepCopyInto(out *GitSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitSource.
func (in *GitSource) DeepCopy() *GitSource {
	if in == nil {
		return nil
	}
	out := new(GitSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostAndPortEndpoint) DeepCopyInto(out *HostAndPortEndpoint) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostAndPortEndpoint.
func (in *HostAndPortEndpoint) DeepCopy() *HostAndPortEndpoint {
	if in == nil {
		return nil
	}
	out := new(HostAndPortEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MessagingEndpoint) DeepCopyInto(out *MessagingEndpoint) {
	*out = *in
	out.HostAndPortEndpoint = in.HostAndPortEndpoint
	if in.CACertificate != nil {
		in, out := &in.CACertificate, &out.CACertificate
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MessagingEndpoint.
func (in *MessagingEndpoint) DeepCopy() *MessagingEndpoint {
	if in == nil {
		return nil
	}
	out := new(MessagingEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProducerSpec) DeepCopyInto(out *ProducerSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
	if in.NumberOfThreads != nil {
		in, out := &in.NumberOfThreads, &out.NumberOfThreads
		*out = new(uint32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProducerSpec.
func (in *ProducerSpec) DeepCopy() *ProducerSpec {
	if in == nil {
		return nil
	}
	out := new(ProducerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProducerStatus) DeepCopyInto(out *ProducerStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProducerStatus.
func (in *ProducerStatus) DeepCopy() *ProducerStatus {
	if in == nil {
		return nil
	}
	out := new(ProducerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegistryEndpoint) DeepCopyInto(out *RegistryEndpoint) {
	*out = *in
	out.URLEndpoint = in.URLEndpoint
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegistryEndpoint.
func (in *RegistryEndpoint) DeepCopy() *RegistryEndpoint {
	if in == nil {
		return nil
	}
	out := new(RegistryEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Simulator) DeepCopyInto(out *Simulator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Simulator.
func (in *Simulator) DeepCopy() *Simulator {
	if in == nil {
		return nil
	}
	out := new(Simulator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Simulator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorConsumer) DeepCopyInto(out *SimulatorConsumer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorConsumer.
func (in *SimulatorConsumer) DeepCopy() *SimulatorConsumer {
	if in == nil {
		return nil
	}
	out := new(SimulatorConsumer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SimulatorConsumer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorConsumerList) DeepCopyInto(out *SimulatorConsumerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SimulatorConsumer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorConsumerList.
func (in *SimulatorConsumerList) DeepCopy() *SimulatorConsumerList {
	if in == nil {
		return nil
	}
	out := new(SimulatorConsumerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SimulatorConsumerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorEndpoint) DeepCopyInto(out *SimulatorEndpoint) {
	*out = *in
	in.Messaging.DeepCopyInto(&out.Messaging)
	out.Registry = in.Registry
	out.Adapters = in.Adapters
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorEndpoint.
func (in *SimulatorEndpoint) DeepCopy() *SimulatorEndpoint {
	if in == nil {
		return nil
	}
	out := new(SimulatorEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorList) DeepCopyInto(out *SimulatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Simulator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorList.
func (in *SimulatorList) DeepCopy() *SimulatorList {
	if in == nil {
		return nil
	}
	out := new(SimulatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SimulatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorProducer) DeepCopyInto(out *SimulatorProducer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorProducer.
func (in *SimulatorProducer) DeepCopy() *SimulatorProducer {
	if in == nil {
		return nil
	}
	out := new(SimulatorProducer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SimulatorProducer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorProducerList) DeepCopyInto(out *SimulatorProducerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SimulatorProducer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorProducerList.
func (in *SimulatorProducerList) DeepCopy() *SimulatorProducerList {
	if in == nil {
		return nil
	}
	out := new(SimulatorProducerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SimulatorProducerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorSpec) DeepCopyInto(out *SimulatorSpec) {
	*out = *in
	if in.Builds != nil {
		in, out := &in.Builds, &out.Builds
		*out = make(map[string]Build, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Endpoint.DeepCopyInto(&out.Endpoint)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorSpec.
func (in *SimulatorSpec) DeepCopy() *SimulatorSpec {
	if in == nil {
		return nil
	}
	out := new(SimulatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorStatus) DeepCopyInto(out *SimulatorStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorStatus.
func (in *SimulatorStatus) DeepCopy() *SimulatorStatus {
	if in == nil {
		return nil
	}
	out := new(SimulatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *URLEndpoint) DeepCopyInto(out *URLEndpoint) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new URLEndpoint.
func (in *URLEndpoint) DeepCopy() *URLEndpoint {
	if in == nil {
		return nil
	}
	out := new(URLEndpoint)
	in.DeepCopyInto(out)
	return out
}
