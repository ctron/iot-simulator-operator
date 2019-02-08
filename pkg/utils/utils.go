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
	"fmt"

	"github.com/ctron/iot-simulator-operator/pkg/apis/simulator/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetOwnerReference(owner v1.Object, existingObject runtime.Object, scheme *runtime.Scheme) error {

	// we need to roll our own logic here, since we do not want
	// to set blockOwnerDeletion to TRUE

	ro, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("is not a %T a runtime.Object, cannot call SetControllerReference", owner)
	}

	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		return err
	}

	existingObj := existingObject.(v1.Object)
	var TRUE = true
	ts := existingObj.GetCreationTimestamp()
	if ts.IsZero() {
		existingObj.SetOwnerReferences([]v1.OwnerReference{{
			APIVersion: gvk.GroupVersion().String(),
			Kind:       gvk.Kind,
			Name:       owner.GetName(),
			UID:        owner.GetUID(),
			Controller: &TRUE,
		}})
	}

	return nil
}

func MakeHelmInstanceName(obj v1alpha1.SimulatorComponent) string {
	if obj.GetCommonSpec().Simulator == "" {
		return "iot-simulator"
	} else {
		return obj.GetCommonSpec().Simulator + "-iot-simulator"
	}
}

func DeploymentConfigName(prefix string, obj metav1.Object) string {
	return "dc-" + prefix + "-" + obj.GetName()
}
