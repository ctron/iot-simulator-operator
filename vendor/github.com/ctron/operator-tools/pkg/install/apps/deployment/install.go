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

package deployment

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DeploymentMutator func(*appsv1.Deployment) (reconcile.Result, error)
type DeploymentMutatorSimple func(*appsv1.Deployment) error

func ReconcileDeployment(name string, mutator DeploymentMutator, mixin install.MixIn) recon.Processor {

	obj := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*appsv1.Deployment))
	}, mixin)

}

func SimpleDeployment(mutator DeploymentMutatorSimple) DeploymentMutator {
	return func(config *appsv1.Deployment) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileDeploymentSimple(name string, mutator DeploymentMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileDeployment(name, SimpleDeployment(mutator), mixin)
}
