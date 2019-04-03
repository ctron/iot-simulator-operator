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

package images

import (
	"os"
	"strings"

	"github.com/ctron/iot-simulator-operator/pkg/apis/simulator/v1alpha1"
)

var SimulatorImage string
var ConsoleImage string

var defaultGitPrefix string
var defaultGitRef string

func init() {
	base := os.Getenv("IMAGE_BASE")

	if base == "" {
		base = "docker.io/ctron"

	}

	tag := os.Getenv("IMAGE_TAG")
	if tag == "" {
		tag = ":latest"
	} else if !strings.HasPrefix(tag, ":") {
		tag = ":" + tag
	}

	SimulatorImage = base + "/iot-hono-simulator" + tag
	ConsoleImage = base + "/iot-simulator-console" + tag

	// git

	defaultGitPrefix = os.Getenv("DEFAULT_GIT_PREFIX")
	if defaultGitPrefix == "" {
		defaultGitPrefix = "https://github.com/ctron"
	}
	defaultGitRef = os.Getenv("DEFAULT_GIT_REF")
	if defaultGitRef == "" {
		defaultGitRef = "master"
	}
}

func defaultBuildSource(repo string) (string, string) {
	return buildSource(repo, "", "")
}

func buildSource(repo string, uri string, ref string) (string, string) {
	if uri == "" {
		uri = defaultGitPrefix + "/" + repo
	}
	if ref == "" {
		ref = defaultGitRef
	}
	return uri, ref
}

func EvalBuildSource(simulator *v1alpha1.Simulator, repo string) (string, string) {

	if simulator.Spec.Builds == nil {
		return defaultBuildSource(repo)
	}

	b, ok := simulator.Spec.Builds[repo]
	if !ok {
		return defaultBuildSource(repo)
	}

	return buildSource(repo, b.Git.URI, b.Git.Reference)

}
