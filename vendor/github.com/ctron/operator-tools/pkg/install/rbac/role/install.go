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

package role

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type RoleMutator func(*rbacv1.Role) (reconcile.Result, error)
type RoleMutatorSimple func(*rbacv1.Role) error

func ReconcileRole(name string, mutator RoleMutator, mixin install.MixIn) recon.Processor {

	obj := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*rbacv1.Role))
	}, mixin)

}

func SimpleRole(mutator RoleMutatorSimple) RoleMutator {
	return func(config *rbacv1.Role) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileRoleSimple(name string, mutator RoleMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileRole(name, SimpleRole(mutator), mixin)
}
