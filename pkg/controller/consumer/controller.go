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

	mycontroller "github.com/ctron/iot-simulator-operator/pkg/controller"
	"github.com/ctron/iot-simulator-operator/pkg/images"

	"github.com/ctron/operator-tools/pkg/install/openshift"

	"github.com/ctron/iot-simulator-operator/pkg/controller/common"

	"k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ctron/iot-simulator-operator/pkg/utils"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kappsv1 "k8s.io/api/apps/v1"
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

	if err := mycontroller.WatchAll(c, &simv1alpha1.SimulatorConsumer{}); err != nil {
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

	if openshift.IsOpenshift() {
		err = r.reconcileDeploymentConfig(request, instance)
	} else {
		err = r.reconcileDeployment(request, instance)
	}

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
			Name:      utils.DeploymentConfigName("con", instance),
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

	existing.ObjectMeta.Labels["app"] = utils.MakeInstanceName(consumer)
	existing.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("con", existing)
	existing.ObjectMeta.Labels["metrics"] = "iot-simulator"
	existing.ObjectMeta.Labels["iot.simulator"] = consumer.Spec.Simulator

	existing.Spec.Ports = []corev1.ServicePort{
		{
			Name: "metrics",
			Port: 8081, TargetPort: intstr.FromInt(8081),
		},
	}
	existing.Spec.Selector = map[string]string{
		"app":              utils.MakeInstanceName(consumer),
		"deploymentconfig": existing.Name,
	}

}

func (r *ReconcileConsumer) reconcileDeploymentConfig(request reconcile.Request, instance *simv1alpha1.SimulatorConsumer) error {
	dc := appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DeploymentConfigName("con", instance),
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

func (r *ReconcileConsumer) reconcileDeployment(request reconcile.Request, instance *simv1alpha1.SimulatorConsumer) error {
	dc := appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DeploymentConfigName("con", instance),
			Namespace: request.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, &dc, func(existingObject runtime.Object) error {

		if err := utils.SetOwnerReference(instance, existingObject, r.scheme); err != nil {
			return err
		}

		existing := existingObject.(*kappsv1.Deployment)
		r.configureDeployment(instance, existing)

		return nil
	})

	return err
}

func (r *ReconcileConsumer) applyConsumerPodSpec(consumer *simv1alpha1.SimulatorConsumer, obj metav1.Object, pod *corev1.PodTemplateSpec) {

	labels := obj.GetLabels()

	if labels == nil {
		labels = map[string]string{}
	}

	messageType := consumer.Spec.Type
	if messageType == "" {
		messageType = "telemetry"
	}

	labels["app"] = utils.MakeInstanceName(consumer)
	labels["deploymentconfig"] = utils.DeploymentConfigName("con", obj)
	labels["iot.simulator.tenant"] = consumer.Spec.Tenant
	labels["iot.simulator"] = consumer.Spec.Simulator
	labels["iot.simulator.app"] = "consumer"
	labels["iot.simulator.message.type"] = messageType

	obj.SetLabels(labels)

	// template

	simulatorName := consumer.Spec.Simulator

	if pod.ObjectMeta.Labels == nil {
		pod.ObjectMeta.Labels = make(map[string]string)
	}

	pod.ObjectMeta.Labels["app"] = labels["app"]
	pod.ObjectMeta.Labels["deploymentconfig"] = labels["deploymentconfig"]
	pod.ObjectMeta.Labels["iot.simulator.tenant"] = consumer.Spec.Tenant
	pod.ObjectMeta.Labels["iot.simulator"] = consumer.Spec.Simulator

	// container

	if len(pod.Spec.Containers) != 1 {
		pod.Spec.Containers = make([]corev1.Container, 1)
	}

	pod.Spec.Containers[0].Name = "consumer"
	pod.Spec.Containers[0].Command = []string{"java", "-Dvertx.cacheDirBase=/tmp", "-Dvertx.logger-delegate-factory-class-name=io.vertx.core.logging.SLF4JLogDelegateFactory", "-jar", "/build/simulator-consumer/target/simulator-consumer-app.jar"}
	pod.Spec.Containers[0].Env = []v1.EnvVar{
		{Name: "CONSUMING", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.labels['iot.simulator.message.type']"}}},
		{Name: "HONO_TRUSTED_CERTS", Value: "/etc/secrets/messaging.ca.crt"},
		{Name: "HONO_INITIAL_CREDITS", Value: "100"},
		{Name: "HONO_TENANT", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.labels['iot.simulator.tenant']"}}},

		{Name: "HONO_USER", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "endpoint.username"}}},
		{Name: "HONO_PASSWORD", ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "endpoint.password"}}},
		{Name: "MESSAGING_SERVICE_HOST", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "endpoint.host"}}},
		{Name: "MESSAGING_SERVICE_PORT_AMQP", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "endpoint.port"}}},
	}
	pod.Spec.Containers[0].Ports = []v1.ContainerPort{
		{
			ContainerPort: 8081,
			Name:          "metrics",
			Protocol:      corev1.ProtocolTCP,
		},
	}

	// volumes

	if len(pod.Spec.Containers[0].VolumeMounts) != 1 {
		pod.Spec.Containers[0].VolumeMounts = make([]corev1.VolumeMount, 1)
	}
	pod.Spec.Containers[0].VolumeMounts[0].Name = "secrets"
	pod.Spec.Containers[0].VolumeMounts[0].MountPath = "/etc/secrets"

	if len(pod.Spec.Volumes) != 1 {
		pod.Spec.Volumes = make([]corev1.Volume, 1)
	}
	pod.Spec.Volumes[0].Name = "secrets"
	if pod.Spec.Volumes[0].Secret == nil {
		pod.Spec.Volumes[0].Secret = &corev1.SecretVolumeSource{}
	}
	pod.Spec.Volumes[0].Secret.SecretName = consumer.Spec.Simulator

	// health checks

	pod.Spec.Containers[0].LivenessProbe = common.ApplyProbe(pod.Spec.Containers[0].LivenessProbe)
	pod.Spec.Containers[0].ReadinessProbe = common.ApplyProbe(pod.Spec.Containers[0].ReadinessProbe)

	// limits

	pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: *resource.NewQuantity(1024*1024*1024 /* 1024Mi */, resource.BinarySI),
		},
	}

}

func (r *ReconcileConsumer) configureDeployment(consumer *simv1alpha1.SimulatorConsumer, existing *kappsv1.Deployment) {

	r.applyConsumerPodSpec(consumer, existing, &existing.Spec.Template)

	existing.Spec.Replicas = &consumer.Spec.Replicas

	if existing.Spec.Selector == nil {
		existing.Spec.Selector = &metav1.LabelSelector{}
	}

	existing.Spec.Selector.MatchLabels = map[string]string{
		"app":              existing.Labels["app"],
		"deploymentconfig": existing.Labels["deploymentconfig"],
	}

	existing.Spec.Template.Spec.Containers[0].Image = images.SimulatorImage

}

func (r *ReconcileConsumer) configureDeploymentConfig(consumer *simv1alpha1.SimulatorConsumer, existing *appsv1.DeploymentConfig) {

	if existing.Spec.Template == nil {
		existing.Spec.Template = &v1.PodTemplateSpec{}
	}

	r.applyConsumerPodSpec(consumer, existing, existing.Spec.Template)

	existing.Spec.Replicas = consumer.Spec.Replicas
	existing.Spec.Selector = map[string]string{
		"app":              existing.Labels["app"],
		"deploymentconfig": existing.Labels["deploymentconfig"],
	}

	existing.Spec.Strategy.Type = appsv1.DeploymentStrategyTypeRolling

	// triggers

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
	existing.Spec.Triggers[1].ImageChangeParams.From.Kind = "ImageStreamTag"
	existing.Spec.Triggers[1].ImageChangeParams.From.Name = utils.MakeInstanceName(consumer) + "-parent:latest"

}
