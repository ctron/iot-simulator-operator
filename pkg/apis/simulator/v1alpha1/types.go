/*******************************************************************************
 * Copyright (c) 2018, 2019 Red Hat Inc
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Simulator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SimulatorSpec   `json:"spec"`
	Status SimulatorStatus `json:"status"`
}

type SimulatorSpec struct {
	Builds   map[string]Build  `json:"builds,omitempty"`
	Endpoint SimulatorEndpoint `json:"endpoint"`
}

type SimulatorStatus struct {
}

type Build struct {
	Git GitSource `json:"git"`
}

type GitSource struct {
	URI       string `json:"uri"`
	Reference string `json:"ref,omitempty"`
}

type SimulatorEndpoint struct {
	Messaging MessagingEndpoint `json:"messaging"`
	Registry  URLEndpoint       `json:"registry"`
	Adapters  AdapterEndpoints  `json:"adapters"`
}

type HostAndPortEndpoint struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type MessagingEndpoint struct {
	HostAndPortEndpoint `json:",inline"`

	User          string `json:"user"`
	Password      string `json:"password"`
	CACertificate []byte `json:"caCertificate"`
}

type AdapterEndpoints struct {
	HTTP URLEndpoint         `json:"http"`
	MQTT HostAndPortEndpoint `json:"mqtt"`
}

type URLEndpoint struct {
	URL string `json:"url"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SimulatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Simulator `json:"items"`
}

type CommonSpec struct {
	Simulator string `json:"simulator"`
	Tenant    string `json:"tenant"`
	Type      string `json:"type"`

	Replicas *int32 `json:"replicas,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type SimulatorConsumer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConsumerSpec   `json:"spec"`
	Status ConsumerStatus `json:"status"`
}

type ConsumerSpec struct {
	CommonSpec
}

type ConsumerStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true

type SimulatorProducer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProducerSpec   `json:"spec"`
	Status ProducerStatus `json:"status"`
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

type ProducerStatus struct {
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
		&Simulator{},
		&SimulatorList{},

		&SimulatorConsumer{},
		&SimulatorConsumerList{},

		&SimulatorProducer{},
		&SimulatorProducerList{},
	)
}
