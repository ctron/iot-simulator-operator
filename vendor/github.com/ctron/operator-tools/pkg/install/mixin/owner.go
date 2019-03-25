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

package mixin

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	"github.com/ctron/operator-tools/pkg/utils/owner"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ObjectOwnerReference(owner owner.OwnerMutator) install.MixIn {
	return func(ctx recon.InstallContext, object metav1.Object) error {
		return owner(ctx, object)
	}
}

func ControllerOwner(ownerObject metav1.Object) install.MixIn {
	return ObjectOwnerReference(owner.ControllerOwner(ownerObject))
}

func ObjectOwner(ownerObject metav1.Object) install.MixIn {
	return ObjectOwnerReference(owner.ObjectOwner(ownerObject))
}
