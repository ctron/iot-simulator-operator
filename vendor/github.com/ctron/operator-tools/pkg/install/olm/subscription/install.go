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

package subscription

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	subv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SubscriptionMutator func(*subv1.Subscription) (reconcile.Result, error)
type SubscriptionMutatorSimple func(*subv1.Subscription) error

func ReconcileSubscription(name string, mutator SubscriptionMutator, mixin install.MixIn) recon.Processor {

	obj := subv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*subv1.Subscription))
	}, mixin)

}

func SimpleSubscription(mutator SubscriptionMutatorSimple) SubscriptionMutator {
	return func(config *subv1.Subscription) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileSubscriptionSimple(name string, mutator SubscriptionMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileSubscription(name, SimpleSubscription(mutator), mixin)
}
