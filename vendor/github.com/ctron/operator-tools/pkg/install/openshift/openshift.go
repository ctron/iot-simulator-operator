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

package openshift

import (
	"log"
	"os"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client/config"

	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
)

var (
	openshift *bool
)

func IsOpenshift() bool {
	if openshift == nil {
		b := detectOpenshift()
		openshift = &b
	}
	return *openshift
}

func detectOpenshift() bool {

	value := os.Getenv("USE_OPENSHIFT")
	if value != "" {
		v := strings.ToLower(value)
		log.Printf("USE_OPENSHIFT = %s", v)
		return v == "true"
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("Error getting config: %v", err)
		return false
	}

	routeClient, err := routev1.NewForConfig(cfg)
	if err != nil {
		log.Printf("Failed to get routeClient: %v", err)
		return false
	}

	body, err := routeClient.RESTClient().Get().DoRaw()

	log.Printf("Request error: %v", err)
	log.Printf("Body: %v", string(body))

	return err == nil
}
