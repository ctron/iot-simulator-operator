/*******************************************************************************
 * Copyright (c) 2019 Red Hat Inc
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

package utils

import (
	"github.com/ctron/iot-simulator-operator/pkg/apis/simulator/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func SetOwnerReference(instance *v1alpha1.SimulatorConsumer, existingObject runtime.Object) {

	// we need to roll our own logic here, since we do not want
	// to set blockOwnerDeletion to TRUE

	existingObj := existingObject.(v1.Object)
	var TRUE = true
	ts := existingObj.GetCreationTimestamp()
	if ts.IsZero() {
		existingObj.SetOwnerReferences([]v1.OwnerReference{{
			APIVersion: instance.APIVersion,
			Kind:       instance.Kind,
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
			Controller: &TRUE,
		}})
	}

}

func MakeHelmInstanceName(consumer *v1alpha1.SimulatorConsumer) string {
	if consumer.Spec.Simulator == "" {
		return "iot-simulator"
	} else {
		return consumer.Spec.Simulator + "-iot-simulator"
	}
}
