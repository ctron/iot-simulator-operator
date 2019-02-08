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

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ctron/iot-simulator-operator/pkg/utils"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"

	simv1alpha1 "github.com/ctron/iot-simulator-operator/pkg/apis/simulator/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
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

	err = c.Watch(&source.Kind{Type: &simv1alpha1.SimulatorConsumer{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
		IsController: true, OwnerType: &simv1alpha1.SimulatorConsumer{},
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
	instance := &simv1alpha1.SimulatorConsumer{}
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

	err = r.reconcileDeploymentConfig(request, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileService(request, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileConsumer) reconcileService(request reconcile.Request, instance *simv1alpha1.SimulatorConsumer) error {

	svc := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sim-consumer-" + request.Name,
			Namespace: request.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, &svc, func(existingObject runtime.Object) error {

		utils.SetOwnerReference(instance, existingObject)

		existing := existingObject.(*v1.Service)
		r.configureService(instance, existing)

		return nil
	})

	return err
}

func (r *ReconcileConsumer) configureService(consumer *simv1alpha1.SimulatorConsumer, existing *v1.Service) {

	existing.Spec = v1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{Name: "metrics", Port: 8081, TargetPort: intstr.FromInt(8081)},
		},
		Selector: map[string]string{
			"app":              utils.MakeHelmInstanceName(consumer),
			"deploymentconfig": "dc-" + existing.Name,
			"metrics":          utils.MakeHelmInstanceName(consumer),
		},
	}

}

func (r *ReconcileConsumer) reconcileDeploymentConfig(request reconcile.Request, instance *simv1alpha1.SimulatorConsumer) error {
	dc := appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sim-consumer-" + request.Name,
			Namespace: request.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, &dc, func(existingObject runtime.Object) error {

		utils.SetOwnerReference(instance, existingObject)

		existing := existingObject.(*appsv1.DeploymentConfig)
		r.configureDeploymentConfig(instance, existing)

		return nil
	})

	return err
}

func (r *ReconcileConsumer) configureDeploymentConfig(consumer *simv1alpha1.SimulatorConsumer, existing *appsv1.DeploymentConfig) {

	sec := consumer.Spec.EndpointSecret

	existing.Spec = appsv1.DeploymentConfigSpec{
		Replicas: 1,
		Selector: map[string]string{
			"app":              utils.MakeHelmInstanceName(consumer),
			"deploymentconfig": "dc-" + existing.Name,
		},
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.DeploymentStrategyTypeRolling,
		},
		Template: &corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app":                  utils.MakeHelmInstanceName(consumer),
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
							{Name: "HONO_TRUSTED_CERTS", Value: "/etc/secrets/ca.crt"},
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
			{Type: appsv1.DeploymentTriggerOnConfigChange},
			{Type: appsv1.DeploymentTriggerOnImageChange, ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
				Automatic:      true,
				ContainerNames: []string{"consumer"},
				From: v1.ObjectReference{
					Kind: "ImageStreamTag",
					Name: utils.MakeHelmInstanceName(consumer) + ":latest",
				},
			}},
		},
	}
}
