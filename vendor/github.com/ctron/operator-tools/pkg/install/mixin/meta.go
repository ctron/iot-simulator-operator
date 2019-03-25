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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AppendLabelMixin(key string, value string) install.MixIn {
	return AppendLabelsMixin(map[string]string{key: value})
}

func AppendLabelsMixin(add map[string]string) install.MixIn {
	return func(ctx recon.InstallContext, object metav1.Object) error {
		m := object.GetLabels()

		if m == nil {
			m = make(map[string]string)
		}

		for k, v := range add {
			m[k] = v
		}

		object.SetLabels(m)

		return nil
	}
}

func AppendAnnotationMixin(key string, value string) install.MixIn {
	return AppendAnnotationsMixin(map[string]string{key: value})
}

func AppendAnnotationsMixin(add map[string]string) install.MixIn {
	return func(ctx recon.InstallContext, object metav1.Object) error {
		m := object.GetAnnotations()

		if m == nil {
			m = make(map[string]string)
		}

		for k, v := range add {
			m[k] = v
		}

		object.SetAnnotations(m)

		return nil
	}
}
