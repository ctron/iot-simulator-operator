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

package utils

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func SetOwnerReference(owner, object metav1.Object, scheme *runtime.Scheme, blockOwnerDeletion bool, controller bool) error {

	ro, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("'owner' must be of type 'runtime.Object' to call 'SetOwnerReference', but is: %T", owner)
	}

	// split GVK

	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		return err
	}

	// new entry

	ref := metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: &blockOwnerDeletion,
		Controller:         &controller,
	}

	// find existing entry

	refs := object.GetOwnerReferences()
	idx := -1
	for i, r := range refs {
		if EqualOwnerReference(ref, r) {
			idx = i
		}
	}

	// append or replace

	if idx == -1 {
		refs = append(refs, ref)
	} else {
		refs[idx] = ref
	}

	// update target

	object.SetOwnerReferences(refs)

	// return

	return nil
}

// Test if two references refer to the same owner
func EqualOwnerReference(o1, o2 metav1.OwnerReference) bool {

	if o1.Name != o2.Name || o1.Kind != o2.Kind {
		return false
	}

	gv1, err := schema.ParseGroupVersion(o1.APIVersion)
	if err != nil {
		return false
	}

	gv2, err := schema.ParseGroupVersion(o2.APIVersion)
	if err != nil {
		return false
	}

	return gv1 == gv2
}
