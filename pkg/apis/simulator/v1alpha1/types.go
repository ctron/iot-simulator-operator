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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CommonSpec struct {
	Simulator string `json:"simulator"`
	Tenant    string `json:"tenant"`
	Type      string `json:"type"`

	Replicas int32 `json:"replicas"`

	EndpointSettings string `json:"endpointSettings"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type SimulatorConsumer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConsumerSpec `json:"spec,omitempty"`
}

type ConsumerSpec struct {
	CommonSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true

type SimulatorProducer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProducerSpec `json:"spec,omitempty"`
}

type Protocol string

const ProtocolHttp Protocol = "http"
const ProtocolMqtt Protocol = "mqtt"

type ProducerSpec struct {
	CommonSpec

	Protocol Protocol `json:"protocol"`

	NumberOfDevices uint32  `json:"numberOfDevices"`
	NumberOfThreads *uint32 `json:"numberOfThreads,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SimulatorConsumerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SimulatorConsumer `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SimulatorProducerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SimulatorProducer `json:"items"`
}

// init

func init() {
	SchemeBuilder.Register(
		&SimulatorConsumer{},
		&SimulatorConsumerList{},

		&SimulatorProducer{},
		&SimulatorProducerList{},
	)
}
