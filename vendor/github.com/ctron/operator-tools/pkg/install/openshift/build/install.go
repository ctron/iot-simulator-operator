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

package build

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	buildv1 "github.com/openshift/api/build/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type BuildConfigMutator func(*buildv1.BuildConfig) (reconcile.Result, error)
type BuildConfigMutatorSimple func(*buildv1.BuildConfig) error

func ReconcileBuildConfig(name string, mutator BuildConfigMutator, mixin install.MixIn) recon.Processor {

	obj := buildv1.BuildConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*buildv1.BuildConfig))
	}, mixin)

}

func SimpleBuildConfig(mutator BuildConfigMutatorSimple) BuildConfigMutator {
	return func(config *buildv1.BuildConfig) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileBuildConfigSimple(name string, mutator BuildConfigMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileBuildConfig(name, SimpleBuildConfig(mutator), mixin)
}
