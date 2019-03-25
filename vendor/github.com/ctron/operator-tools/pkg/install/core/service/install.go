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

package service

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ServiceMutator func(*corev1.Service) (reconcile.Result, error)
type ServiceMutatorSimple func(*corev1.Service) error

func ReconcileService(name string, mutator ServiceMutator, mixin install.MixIn) recon.Processor {

	obj := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*corev1.Service))
	}, mixin)

}

func SimpleService(mutator ServiceMutatorSimple) ServiceMutator {
	return func(config *corev1.Service) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileServiceSimple(name string, mutator ServiceMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileService(name, SimpleService(mutator), mixin)
}
