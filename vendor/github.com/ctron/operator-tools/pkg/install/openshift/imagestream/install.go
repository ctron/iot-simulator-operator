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

package imagestream

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	imgv1 "github.com/openshift/api/image/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ImageStreamMutator func(*imgv1.ImageStream) (reconcile.Result, error)
type ImageStreamMutatorSimple func(*imgv1.ImageStream) error

func ReconcileImageStream(name string, mutator ImageStreamMutator, mixin install.MixIn) recon.Processor {

	obj := imgv1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*imgv1.ImageStream))
	}, mixin)

}

func SimpleImageStream(mutator ImageStreamMutatorSimple) ImageStreamMutator {
	return func(config *imgv1.ImageStream) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileImageStreamSimple(name string, mutator ImageStreamMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileImageStream(name, SimpleImageStream(mutator), mixin)
}
