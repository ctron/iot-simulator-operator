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
	"github.com/ctron/operator-tools/pkg/recon"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func DeleteOperation(object runtime.Object) recon.Processor {
	return recon.Simple(func(ctx recon.InstallContext) error {
		err := ctx.GetClient().Delete(ctx.GetContext(), object)
		if err == nil || errors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	})
}
