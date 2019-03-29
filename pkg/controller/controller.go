/*******************************************************************************
 * Copyright (c) 2018 Red Hat Inc
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

package controller

import (
	"github.com/ctron/operator-tools/pkg/install/openshift"
	appsv1 "github.com/openshift/api/apps/v1"
	buildv1 "github.com/openshift/api/build/v1"
	imgv1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	kappsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}

func WatchAll(c controller.Controller, ownerType runtime.Object) error {

	owners := []runtime.Object{
		&corev1.Service{},
		&rbacv1.Role{},
		&rbacv1.RoleBinding{},
	}

	if openshift.IsOpenshift() {
		owners = append(owners, []runtime.Object{
			&appsv1.DeploymentConfig{},
			&buildv1.BuildConfig{},
			&routev1.Route{},
			&imgv1.ImageStream{},
		}...)
	} else {
		owners = append(owners, []runtime.Object{
			&kappsv1.Deployment{},
		}...)
	}

	for _, i := range owners {
		if err := c.Watch(&source.Kind{Type: i}, &handler.EnqueueRequestForOwner{
			IsController: true, OwnerType: ownerType,
		}); err != nil {
			return err
		}
	}

	return nil
}
