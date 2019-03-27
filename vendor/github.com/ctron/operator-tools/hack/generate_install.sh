#!/usr/bin/env bash

set -e
set -o errexit
set -o nounset
set -o pipefail

SCRIPTPATH="$(cd "$(dirname "$0")" && pwd -P)"
BASE=${SCRIPTPATH}/../

function generate() {
    PREFIX="$1"
    PKG="$2"
    KIND="$3"
    ALIAS="$4"
    IMPORT="$5"

    FILE="$BASE/pkg/install/$PREFIX/$PKG/install.go"
    cat $SCRIPTPATH/header.txt > "$FILE"
    (
    cat <<__EOF__
package $PKG

import (
	"github.com/ctron/operator-tools/pkg/install"
	"github.com/ctron/operator-tools/pkg/recon"
	$ALIAS "$IMPORT"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ${KIND}Mutator func(*${ALIAS}.${KIND}) (reconcile.Result, error)
type ${KIND}MutatorSimple func(*${ALIAS}.${KIND}) error

func Reconcile${KIND}(name string, mutator ${KIND}Mutator, mixin install.MixIn) recon.Processor {

	obj := ${ALIAS}.${KIND}{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	return install.ReconcileObject(&obj, func(existingObject runtime.Object) (result reconcile.Result, e error) {
		return mutator(existingObject.(*${ALIAS}.${KIND}))
	}, mixin)

}

func Simple${KIND}(mutator ${KIND}MutatorSimple) ${KIND}Mutator {
	return func(config *${ALIAS}.${KIND}) (reconcile.Result, error) {
		return reconcile.Result{}, mutator(config)
	}
}

func Reconcile${KIND}Simple(name string, mutator ${KIND}MutatorSimple, mixin install.MixIn) recon.Processor {
	return Reconcile${KIND}(name, Simple${KIND}(mutator), mixin)
}

__EOF__
    ) >> "$FILE"

    go fmt "$FILE"
}

generate "core" serviceaccount ServiceAccount corev1 "k8s.io/api/core/v1"
generate "core" service Service corev1 "k8s.io/api/core/v1"
generate "core" configmap ConfigMap corev1 "k8s.io/api/core/v1"
generate "core" secret Secret corev1 "k8s.io/api/core/v1"

generate "rbac" role Role rbacv1 "k8s.io/api/rbac/v1"
generate "rbac" rolebinding RoleBinding rbacv1 "k8s.io/api/rbac/v1"

generate "openshift" build BuildConfig buildv1 "github.com/openshift/api/build/v1"
generate "openshift" imagestream ImageStream imgv1 "github.com/openshift/api/image/v1"
generate "openshift" dc DeploymentConfig appsv1 "github.com/openshift/api/apps/v1"
generate "openshift" route Route routev1 "github.com/openshift/api/route/v1"

generate "olm" subscription Subscription subv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
