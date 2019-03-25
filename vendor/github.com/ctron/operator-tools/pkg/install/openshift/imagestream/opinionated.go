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

package imagestream

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	imgv1 "github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func EmptyImageStream(name string, mixin install.MixIn) recon.Processor {
	return ReconcileImageStream(name, func(stream *imgv1.ImageStream) (reconcile.Result, error) {
		return reconcile.Result{}, nil
	}, mixin)
}

func DockerImageStream(name string, tag string, imageName string, mixin install.MixIn) recon.Processor {
	return ReconcileImageStream(name, func(stream *imgv1.ImageStream) (reconcile.Result, error) {

		stream.Spec.LookupPolicy.Local = false

		return EnsureTag(stream, tag, DockerImageTagReference(imageName))
	}, mixin)
}

type TagReferenceMutator func(reference *imgv1.TagReference) (reconcile.Result, error)

func DockerImageTagReference(imageName string) TagReferenceMutator {

	return func(reference *imgv1.TagReference) (reconcile.Result, error) {

		reference.From = &corev1.ObjectReference{
			Kind: "DockerImage",
			Name: imageName,
		}

		reference.ImportPolicy.Scheduled = true
		reference.ReferencePolicy.Type = imgv1.SourceTagReferencePolicy

		return reconcile.Result{}, nil
	}

}

func EnsureTag(stream *imgv1.ImageStream, name string, mutator TagReferenceMutator) (reconcile.Result, error) {

	if stream.Spec.Tags == nil {
		stream.Spec.Tags = make([]imgv1.TagReference, 0)
	}

	for i, v := range stream.Spec.Tags {
		if v.Name == name {
			return mutator(&stream.Spec.Tags[i])
		}
	}

	v := &imgv1.TagReference{
		Name: name,
	}

	result, err := mutator(v)

	if err == nil {
		stream.Spec.Tags = append(stream.Spec.Tags, *v)
	}

	return result, err
}
