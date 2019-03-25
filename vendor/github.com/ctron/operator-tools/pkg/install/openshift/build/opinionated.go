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

package build

import (
	buildv1 "github.com/openshift/api/build/v1"
	"k8s.io/api/core/v1"
)

// Set git as source
func SetGitSource(config *buildv1.BuildConfig, uri string, ref string) {

	config.Spec.Source.Type = buildv1.BuildSourceGit
	if config.Spec.Source.Git == nil {
		config.Spec.Source.Git = &buildv1.GitBuildSource{}
	}
	config.Spec.Source.Git.URI = uri
	config.Spec.Source.Git.Ref = ref

}

// Set docker as build strategy
func SetDockerStrategyFromImageStream(config *buildv1.BuildConfig, imageStreamTag string) {

	config.Spec.Strategy.Type = buildv1.DockerBuildStrategyType

	if config.Spec.Strategy.DockerStrategy == nil {
		config.Spec.Strategy.DockerStrategy = &buildv1.DockerBuildStrategy{}
	}

	if config.Spec.Strategy.DockerStrategy.From == nil {
		config.Spec.Strategy.DockerStrategy.From = &v1.ObjectReference{}
	}

	config.Spec.Strategy.DockerStrategy.From.Kind = "ImageStreamTag"
	config.Spec.Strategy.DockerStrategy.From.Name = imageStreamTag
}

// Set image stream tag as output
func SetOutputImageStream(config *buildv1.BuildConfig, imageStreamTag string) {

	if config.Spec.Output.To == nil {
		config.Spec.Output.To = &v1.ObjectReference{}
	}

	config.Spec.Output.To.Kind = "ImageStreamTag"
	config.Spec.Output.To.Name = imageStreamTag

}

func EnableDefaultTriggers(config *buildv1.BuildConfig) {

	if config.Spec.Triggers == nil {
		config.Spec.Triggers = make([]buildv1.BuildTriggerPolicy, 0)
	}

	needImageChange := true
	needConfigChange := true

	for _, t := range config.Spec.Triggers {

		switch t.Type {
		case buildv1.ImageChangeBuildTriggerType:
			needImageChange = false
		case buildv1.ConfigChangeBuildTriggerType:
			needConfigChange = false
		}

	}

	if needConfigChange {
		config.Spec.Triggers = append(config.Spec.Triggers, buildv1.BuildTriggerPolicy{
			Type: buildv1.ConfigChangeBuildTriggerType,
		})
	}

	if needImageChange {
		config.Spec.Triggers = append(config.Spec.Triggers, buildv1.BuildTriggerPolicy{
			Type: buildv1.ImageChangeBuildTriggerType,
		})
	}

}
