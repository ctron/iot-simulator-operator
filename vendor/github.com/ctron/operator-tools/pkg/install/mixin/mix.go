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
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Allow to combine multiple mixins
func Mix(mixins ...install.MixIn) install.MixIn {
	return func(ctx recon.InstallContext, object v1.Object) error {

		for _, m := range mixins {
			if err := m(ctx, object); err != nil {
				return err
			}
		}

		return nil

	}
}
