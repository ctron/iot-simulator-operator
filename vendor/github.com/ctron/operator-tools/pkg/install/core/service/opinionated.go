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

package service

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	corev1 "k8s.io/api/core/v1"
)

func Service(name string, serviceLabels map[string]string, mutator ServiceMutatorSimple, mixin install.MixIn) recon.Processor {
	return ReconcileServiceSimple(name, func(service *corev1.Service) error {

		if service.ObjectMeta.Labels == nil {
			service.ObjectMeta.Labels = map[string]string{}
		}

		for k, v := range serviceLabels {
			service.ObjectMeta.Labels[k] = v
		}

		service.Spec.Selector = serviceLabels

		if mutator != nil {
			return mutator(service)
		} else {
			return nil
		}
	}, mixin)
}
