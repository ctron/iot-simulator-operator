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

package prometheus

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ServiceMonitorMutator func(*promv1.ServiceMonitor) (reconcile.Result, error)
type ServiceMonitorSimpleMutator func(*promv1.ServiceMonitor) error

func ReconcileServiceMonitor(name string, mutator ServiceMonitorMutator, mixin install.MixIn) recon.Processor {

	obj := promv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*promv1.ServiceMonitor))
	}, mixin)

}

func SimpleServiceMonitor(mutator ServiceMonitorSimpleMutator) ServiceMonitorMutator {
	return func(config *promv1.ServiceMonitor) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func ReconcileServiceMonitorSimple(name string, mutator ServiceMonitorSimpleMutator, mixin install.MixIn) recon.Processor {
	return ReconcileServiceMonitor(name, SimpleServiceMonitor(mutator), mixin)
}
