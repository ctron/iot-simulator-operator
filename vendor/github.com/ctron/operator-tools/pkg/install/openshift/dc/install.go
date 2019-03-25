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

package dc

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	appsv1 "github.com/openshift/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DeploymentConfigMutator func(*appsv1.DeploymentConfig) (reconcile.Result, error)
type DeploymentConfigMutatorSimple func(*appsv1.DeploymentConfig) error

func ReconcileDeploymentConfig(name string, mutator DeploymentConfigMutator, mixin install.MixIn) recon.Processor {

	obj := appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*appsv1.DeploymentConfig))
	}, mixin)

}

func SimpleDeploymentConfig(mutator DeploymentConfigMutatorSimple) DeploymentConfigMutator {
	return func(config *appsv1.DeploymentConfig) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileDeploymentConfigSimple(name string, mutator DeploymentConfigMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileDeploymentConfig(name, SimpleDeploymentConfig(mutator), mixin)
}
