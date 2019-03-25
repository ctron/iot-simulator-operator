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

package install

import (
	"fmt"

	"github.com/ctron/operator-tools/pkg/recon"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Mutator func(existingObject runtime.Object) (reconcile.Result, error)
type MixIn func(ctx recon.InstallContext, object metav1.Object) error

func ReconcileObject(obj runtime.Object, mutator Mutator, mixins ...MixIn) recon.Processor {

	return func(ctx recon.InstallContext) (result reconcile.Result, e error) {

		o, ok := obj.(metav1.Object)
		if !ok {
			return reconcile.Result{}, fmt.Errorf("object is not of type v1.Object")
		}

		o.SetNamespace(ctx.GetRequest().Namespace)

		var r reconcile.Result

		_, err := controllerutil.CreateOrUpdate(ctx.GetContext(), ctx.GetClient(), obj, func(existingObject runtime.Object) error {

			obj, ok := existingObject.(metav1.Object)
			if !ok {
				return fmt.Errorf("object is not of type v1.Object")
			}

			var err error
			r, err = mutator(existingObject)
			if err != nil {
				return err
			}

			for _, m := range mixins {
				err = m(ctx, obj)
				if err != nil {
					return err
				}
			}

			return nil
		})

		return r, err

	}

}
