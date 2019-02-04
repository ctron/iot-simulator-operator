/*******************************************************************************
 * Copyright (c) 2018 Red Hat Inc
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

package consumer

import (
	"context"

	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"

	simv1alpha1 "github.com/ctron/iot-simulator-operator/pkg/apis/simulator/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_consumer")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileConsumer{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {

	c, err := controller.New("consumer-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &simv1alpha1.Consumer{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
		IsController: true, OwnerType: &simv1alpha1.Consumer{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileConsumer{}

type ReconcileConsumer struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileConsumer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Consumer")

	// Fetch the Consumer instance
	instance := &simv1alpha1.Consumer{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	pod := appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sim-consumer-" + request.Name,
			Namespace: request.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(context.TODO(), r.client, &pod, func(existingObject runtime.Object) error {
		existing := existingObject.(*appsv1.DeploymentConfig)

		ts := existing.GetCreationTimestamp()
		if ts.IsZero() {
			if err := controllerutil.SetControllerReference(instance, existing, r.scheme); err != nil {
				return err
			}
		}

		r.reconileDeploymentConfig(instance, existing)

		return nil
	})

	return reconcile.Result{}, err

}

func (r *ReconcileConsumer) reconileDeploymentConfig(consumer *simv1alpha1.Consumer, existing *appsv1.DeploymentConfig) {

	sec := "simulator-secrets-" + consumer.Spec.MessagingEndpoint

	existing.Spec = appsv1.DeploymentConfigSpec{
		Replicas: 1,
		Selector: map[string]string{
			"app":              "simulator",
			"deploymentconfig": "dc-" + existing.Name,
		},
		Strategy: appsv1.DeploymentStrategy{
			Type: "rolling",
		},
		Template: &corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app":                  "simulator",
					"deploymentconfig":     "dc-" + existing.Name,
					"iot.simulator.tenant": consumer.Spec.Tenant,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "consumer",
						Command: []string{"java", "-Xmx1024m", "-Dvertx.cacheDirBase=/tmp", "-Dvertx.logger-delegate-factory-class-name=io.vertx.core.logging.SLF4JLogDelegateFactory", "-jar", "/build/simulator-consumer/target/simulator-consumer-app.jar"},
						Env: []v1.EnvVar{
							{Name: "HONO_TRUSTED_CERTS", Value: "/etc/secrets/server-cert.pem"},
							{Name: "HONO_INITIAL_CREDITS", Value: "100"},
							{Name: "HONO_TENANT", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.labels['iot.simulator.tenant']"}}},
							{Name: "HONO_USER", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: sec}, Key: "endpoint.username"}}},
							{Name: "HONO_PASSWORD", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: sec}, Key: "endpoint.password"}}},
							{Name: "MESSAGING_SERVICE_HOST", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: sec}, Key: "endpoint.host"}}},
							{Name: "MESSAGING_SERVICE_PORT_AMQP", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: sec}, Key: "endpoint.port"}}},
						},
						Ports: []v1.ContainerPort{
							{ContainerPort: 8081, Name: "metrics"},
						},
						VolumeMounts: []corev1.VolumeMount{
							{MountPath: "/etc/secrets", Name: "secrets-volume"},
						},
					},
				},
				Volumes: []corev1.Volume{
					{Name: "secrets-volume", VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: sec,
						},
					}},
				},
			},
		},
		Triggers: appsv1.DeploymentTriggerPolicies{
			{Type: "ConfigChange"},
			{Type: "ImageChange", ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
				Automatic:      true,
				ContainerNames: []string{"consumer"},
				From: v1.ObjectReference{
					Kind: "ImageStreamTag",
					Name: "simulator-parent:latest",
				},
			}},
		},
	}
}
