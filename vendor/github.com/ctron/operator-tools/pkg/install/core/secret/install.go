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

package secret

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SecretMutator func(*corev1.Secret) (reconcile.Result, error)
type SecretMutatorSimple func(*corev1.Secret) error

func ReconcileSecret(name string, mutator SecretMutator, mixin install.MixIn) recon.Processor {

	obj := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*corev1.Secret))
	}, mixin)

}

func SimpleSecret(mutator SecretMutatorSimple) SecretMutator {
	return func(config *corev1.Secret) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileSecretSimple(name string, mutator SecretMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileSecret(name, SimpleSecret(mutator), mixin)
}
