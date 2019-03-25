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

package route

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type RouteMutator func(*routev1.Route) (reconcile.Result, error)
type RouteMutatorSimple func(*routev1.Route) error

func ReconcileRoute(name string, mutator RouteMutator, mixin install.MixIn) recon.Processor {

	obj := routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*routev1.Route))
	}, mixin)

}

func SimpleRoute(mutator RouteMutatorSimple) RouteMutator {
	return func(config *routev1.Route) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileRouteSimple(name string, mutator RouteMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileRoute(name, SimpleRoute(mutator), mixin)
}
