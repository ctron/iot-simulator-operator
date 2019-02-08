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

		if err := utils.SetOwnerReference(instance, existingObject, r.scheme); err != nil {
			return err
		}

		existing := existingObject.(*v1.Service)
		r.configureService(instance, existing)

		return nil
	})

	return err
}

func (r *ReconcileConsumer) configureService(consumer *simv1alpha1.SimulatorConsumer, existing *v1.Service) {

	if existing.ObjectMeta.Labels == nil {
		existing.ObjectMeta.Labels = map[string]string{}
	}

	existing.ObjectMeta.Labels["app"] = utils.MakeHelmInstanceName(consumer)
	existing.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("con", existing)
	existing.ObjectMeta.Labels["metrics"] = utils.MakeHelmInstanceName(consumer)
	existing.ObjectMeta.Labels["iot.simulator"] = consumer.Spec.Simulator

	existing.Spec.Ports = []corev1.ServicePort{
		{Name: "metrics", Port: 8081, TargetPort: intstr.FromInt(8081)},
	}
	existing.Spec.Selector = map[string]string{
		"app":              utils.MakeHelmInstanceName(consumer),
		"deploymentconfig": utils.DeploymentConfigName("con", existing),
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

		if err := utils.SetOwnerReference(instance, existingObject, r.scheme); err != nil {
			return err
		}

		existing := existingObject.(*appsv1.DeploymentConfig)
		r.configureDeploymentConfig(instance, existing)

		return nil
	})

	return err
}

func (r *ReconcileConsumer) configureDeploymentConfig(consumer *simv1alpha1.SimulatorConsumer, existing *appsv1.DeploymentConfig) {

	if existing.ObjectMeta.Labels == nil {
		existing.ObjectMeta.Labels = map[string]string{}
	}

	endpointConfigName := consumer.Spec.EndpointConfig
	messageType := consumer.Spec.Type
	if messageType == "" {
		messageType = "telemetry"
	}

	existing.ObjectMeta.Labels["app"] = utils.MakeHelmInstanceName(consumer)
	existing.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("con", existing)
	existing.ObjectMeta.Labels["iot.simulator.tenant"] = consumer.Spec.Tenant
	existing.ObjectMeta.Labels["iot.simulator"] = consumer.Spec.Simulator
	existing.ObjectMeta.Labels["iot.simulator.app"] = "consumer"
	existing.ObjectMeta.Labels["iot.simulator.message.type"] = messageType

	existing.Spec.Replicas = 1
	existing.Spec.Selector = map[string]string{
		"app":              utils.MakeHelmInstanceName(consumer),
		"deploymentconfig": utils.DeploymentConfigName("con", existing),
	}

	existing.Spec.Strategy.Type = appsv1.DeploymentStrategyTypeRolling

	if existing.Spec.Template == nil {
		existing.Spec.Template = &v1.PodTemplateSpec{}
	}
	existing.Spec.Template.ObjectMeta = metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          utils.MakeHelmInstanceName(consumer),
			"deploymentconfig":             utils.DeploymentConfigName("con", existing),
			"iot.simulator.tenant":         consumer.Spec.Tenant,
			"iot.simulator":                consumer.Spec.Simulator,
			"iot.simulator.endpointConfig": consumer.Spec.EndpointConfig,
		},
	}

	if len(existing.Spec.Template.Spec.Containers) != 1 {
		existing.Spec.Template.Spec.Containers = make([]corev1.Container, 1)
	}

	existing.Spec.Template.Spec.Containers[0].Name = "consumer"
	existing.Spec.Template.Spec.Containers[0].Command = []string{"java", "-Xmx1024m", "-Dvertx.cacheDirBase=/tmp", "-Dvertx.logger-delegate-factory-class-name=io.vertx.core.logging.SLF4JLogDelegateFactory", "-jar", "/build/simulator-consumer/target/simulator-consumer-app.jar"}
	existing.Spec.Template.Spec.Containers[0].Env = []v1.EnvVar{
		{Name: "CONSUMING", Value: messageType},
		{Name: "HONO_TRUSTED_CERTS", Value: "/etc/secrets/ca.crt"},
		{Name: "HONO_INITIAL_CREDITS", Value: "100"},
		{Name: "HONO_TENANT", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.labels['iot.simulator.tenant']"}}},
		{Name: "HONO_USER", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: endpointConfigName}, Key: "endpoint.username"}}},
		{Name: "HONO_PASSWORD", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: endpointConfigName}, Key: "endpoint.password"}}},
		{Name: "MESSAGING_SERVICE_HOST", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: endpointConfigName}, Key: "endpoint.host"}}},
		{Name: "MESSAGING_SERVICE_PORT_AMQP", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: endpointConfigName}, Key: "endpoint.port"}}},
	}
	existing.Spec.Template.Spec.Containers[0].Ports = []v1.ContainerPort{
		{ContainerPort: 8081, Name: "metrics"},
	}
	existing.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
		{MountPath: "/etc/secrets", Name: "secrets-volume"},
	}

	if len(existing.Spec.Template.Spec.Volumes) != 1 {
		existing.Spec.Template.Spec.Volumes = make([]corev1.Volume, 1)
	}

	existing.Spec.Template.Spec.Volumes[0].Name = "secrets-volume"
	if existing.Spec.Template.Spec.Volumes[0].Secret == nil {
		existing.Spec.Template.Spec.Volumes[0].Secret = &corev1.SecretVolumeSource{}
	}
	existing.Spec.Template.Spec.Volumes[0].Secret.SecretName = endpointConfigName

	if len(existing.Spec.Triggers) != 2 {
		existing.Spec.Triggers = make([]appsv1.DeploymentTriggerPolicy, 2)
	}

	existing.Spec.Triggers[0].Type = appsv1.DeploymentTriggerOnConfigChange
	existing.Spec.Triggers[1].Type = appsv1.DeploymentTriggerOnImageChange
	if existing.Spec.Triggers[1].ImageChangeParams == nil {
		existing.Spec.Triggers[1].ImageChangeParams = &appsv1.DeploymentTriggerImageChangeParams{}
	}
	existing.Spec.Triggers[1].ImageChangeParams.Automatic = true
	existing.Spec.Triggers[1].ImageChangeParams.ContainerNames = []string{"consumer"}
	existing.Spec.Triggers[1].ImageChangeParams.From = v1.ObjectReference{
		Kind: "ImageStreamTag",
		Name: utils.MakeHelmInstanceName(consumer) + "-parent:latest",
	}

}
