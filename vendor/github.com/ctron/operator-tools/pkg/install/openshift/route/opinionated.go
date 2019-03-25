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

package route

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type TLSConfigMutator func(config *routev1.TLSConfig)

func ReencryptTLSConfig() TLSConfigMutator {
	return func(config *routev1.TLSConfig) {
		config.Termination = routev1.TLSTerminationReencrypt
		config.InsecureEdgeTerminationPolicy = routev1.InsecureEdgeTerminationPolicyRedirect
	}
}

func PassthroughTLSConfig() TLSConfigMutator {
	return func(config *routev1.TLSConfig) {
		config.Termination = routev1.TLSTerminationPassthrough
	}
}

func ReencryptRoute(name string, serviceName string, targetPort intstr.IntOrString, mixin install.MixIn) recon.Processor {
	return Route(name, serviceName, targetPort, ReencryptTLSConfig(), mixin)
}
func PassthroughRoute(name string, serviceName string, targetPort intstr.IntOrString, mixin install.MixIn) recon.Processor {
	return Route(name, serviceName, targetPort, PassthroughTLSConfig(), mixin)
}

func Route(name string, serviceName string, targetPort intstr.IntOrString, tlsConfig TLSConfigMutator, mixin install.MixIn) recon.Processor {
	return ReconcileRouteSimple(name, func(route *routev1.Route) error {

		route.Spec.Port = &routev1.RoutePort{
			TargetPort: targetPort,
		}

		route.Spec.To.Kind = "Service"
		route.Spec.To.Name = serviceName

		if tlsConfig != nil {
			if route.Spec.TLS == nil {
				route.Spec.TLS = &routev1.TLSConfig{}
			}
			tlsConfig(route.Spec.TLS)
		} else {
			route.Spec.TLS = nil
		}

		return nil
	}, mixin)
}
