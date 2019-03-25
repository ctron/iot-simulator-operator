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

package rolebinding

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type RoleBindingMutator func(*rbacv1.RoleBinding) (reconcile.Result, error)
type RoleBindingMutatorSimple func(*rbacv1.RoleBinding) error

func ReconcileRoleBinding(name string, mutator RoleBindingMutator, mixin install.MixIn) recon.Processor {

	obj := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*rbacv1.RoleBinding))
	}, mixin)

}

func SimpleRoleBinding(mutator RoleBindingMutatorSimple) RoleBindingMutator {
	return func(config *rbacv1.RoleBinding) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileRoleBindingSimple(name string, mutator RoleBindingMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileRoleBinding(name, SimpleRoleBinding(mutator), mixin)
}
