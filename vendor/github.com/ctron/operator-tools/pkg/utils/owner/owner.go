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

package owner

import (
	"github.com/ctron/operator-tools/pkg/recon"
	"github.com/ctron/operator-tools/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OwnerMutator func(ctx recon.InstallContext, obj metav1.Object) error

func NoneOwner() OwnerMutator {
	return nil
}

func ControllerOwner(owner metav1.Object) OwnerMutator {
	return ObjectOwnerReference(owner, true, true)
}

func ObjectOwner(owner metav1.Object) OwnerMutator {
	return ObjectOwnerReference(owner, false, false)
}

func ObjectOwnerReference(owner metav1.Object, blockOwnerDeletion bool, controller bool) OwnerMutator {
	return func(ctx recon.InstallContext, obj metav1.Object) error {
		if err := utils.SetOwnerReference(owner, obj, ctx.GetScheme(), blockOwnerDeletion, controller); err != nil {
			return err
		}
		return nil
	}
}
