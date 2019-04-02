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

package producer

import (
	"context"
	"strconv"

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

var log = logf.Log.WithName("controller_producer")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileProducer{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {

	c, err := controller.New("producer-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &simv1alpha1.SimulatorProducer{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	if err := utils.WatchAll(c, &simv1alpha1.SimulatorProducer{}); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileProducer{}

type ReconcileProducer struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileProducer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Producer")

	// Fetch the Producer instance
	instance := &simv1alpha1.SimulatorProducer{}
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

func (r *ReconcileProducer) reconcileService(request reconcile.Request, instance *simv1alpha1.SimulatorProducer) error {

	svc := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DeploymentConfigName("prod", instance),
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

func (r *ReconcileProducer) reconcileDeploymentConfig(request reconcile.Request, instance *simv1alpha1.SimulatorProducer) error {
	dc := appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DeploymentConfigName("prod", instance),
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

func (r *ReconcileProducer) reconcileDeployment(request reconcile.Request, instance *simv1alpha1.SimulatorProducer) error {
	dc := kappsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DeploymentConfigName("prod", instance),
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

func (r *ReconcileProducer) configureService(producer *simv1alpha1.SimulatorProducer, existing *v1.Service) {

	if existing.ObjectMeta.Labels == nil {
		existing.ObjectMeta.Labels = map[string]string{}
	}

	existing.ObjectMeta.Labels["app"] = utils.MakeInstanceName(producer)
	existing.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("prod", existing)
	existing.ObjectMeta.Labels["metrics"] = "iot-simulator"
	existing.ObjectMeta.Labels["iot.simulator"] = producer.Spec.Simulator

	existing.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "metrics",
			Port:       8081,
			TargetPort: intstr.FromInt(8081),
		},
	}
	existing.Spec.Selector = map[string]string{
		"app":              utils.MakeInstanceName(producer),
		"deploymentconfig": existing.Name,
	}

}

func (r *ReconcileProducer) applyPodSpec(producer *simv1alpha1.SimulatorProducer, obj metav1.Object, pod *corev1.PodTemplateSpec) {

	labels := obj.GetLabels()

	if labels == nil {
		labels = map[string]string{}
	}

	obj.SetLabels(labels)

	simulatorName := producer.Spec.Simulator
	messageType := producer.Spec.Type
	if messageType == "" {
		messageType = "telemetry"
	}
	protocol := producer.Spec.Protocol
	if protocol == "" {
		protocol = simv1alpha1.ProtocolHttp
	}

	labels["app"] = utils.MakeInstanceName(producer)
	labels["deploymentconfig"] = utils.DeploymentConfigName("prod", obj)
	labels["iot.simulator.tenant"] = producer.Spec.Tenant
	labels["iot.simulator"] = producer.Spec.Simulator
	labels["iot.simulator.app"] = "producer"
	labels["iot.simulator.message.type"] = messageType
	labels["iot.simulator.producer.protocol"] = string(protocol)

	// template

	if pod.ObjectMeta.Labels == nil {
		pod.ObjectMeta.Labels = make(map[string]string)
	}

	pod.ObjectMeta.Labels["app"] = utils.MakeInstanceName(producer)
	pod.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("prod", obj)
	pod.ObjectMeta.Labels["iot.simulator.tenant"] = producer.Spec.Tenant
	pod.ObjectMeta.Labels["iot.simulator"] = producer.Spec.Simulator
	pod.ObjectMeta.Labels["iot.simulator.message.type"] = messageType

	// containers

	if len(pod.Spec.Containers) != 1 {
		pod.Spec.Containers = make([]corev1.Container, 1)
	}

	pod.Spec.Containers[0].Name = "producer"
	pod.Spec.Containers[0].Env = []v1.EnvVar{

		{Name: "MESSAGE_TYPE", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.labels['iot.simulator.message.type']"}}},

		{Name: "PERIOD_MS", Value: "1000"},
		{Name: "NUM_DEVICES", Value: strconv.FormatUint(uint64(producer.Spec.NumberOfDevices), 10)},

		{Name: "HONO_TENANT", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.labels['iot.simulator.tenant']"}}},
		{Name: "DEVICE_REGISTRY_URL", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "deviceRegistry.url"}}},
	}

	if producer.Spec.NumberOfThreads != nil {
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, v1.EnvVar{
			Name: "NUM_THREADS", Value: strconv.FormatUint(uint64(*producer.Spec.NumberOfThreads), 10),
		})
	}

	pod.Spec.Containers[0].Ports = []v1.ContainerPort{
		{
			ContainerPort: 8081,
			Name:          "metrics",
			Protocol:      corev1.ProtocolTCP,
		},
	}

	// health checks

	pod.Spec.Containers[0].LivenessProbe = common.ApplyProbe(pod.Spec.Containers[0].LivenessProbe)
	pod.Spec.Containers[0].ReadinessProbe = common.ApplyProbe(pod.Spec.Containers[0].ReadinessProbe)

	// resource limits

	pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: *resource.NewQuantity(1024*1024*1024 /* 1024Mi */, resource.BinarySI),
		},
	}

	// now apply protocol specifics

	switch protocol {
	case simv1alpha1.ProtocolMqtt:
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, mqttVariables(simulatorName)...)
		r.configureMqtt(producer, pod, messageType)
	default:
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, httpVariables(simulatorName)...)
		r.configureHttp(producer, pod, messageType)
	}

}

func (r *ReconcileProducer) configureDeployment(producer *simv1alpha1.SimulatorProducer, existing *kappsv1.Deployment) {

	r.applyPodSpec(producer, existing, &existing.Spec.Template)

	existing.Spec.Replicas = &producer.Spec.Replicas

	if existing.Spec.Selector == nil {
		existing.Spec.Selector = &metav1.LabelSelector{}
	}

	existing.Spec.Selector.MatchLabels = map[string]string{
		"app":              existing.Labels["app"],
		"deploymentconfig": existing.Labels["deploymentconfig"],
	}

	existing.Spec.Template.Spec.Containers[0].Image = images.SimulatorImage
}

func (r *ReconcileProducer) configureDeploymentConfig(producer *simv1alpha1.SimulatorProducer, existing *appsv1.DeploymentConfig) {

	if existing.Spec.Template == nil {
		existing.Spec.Template = &v1.PodTemplateSpec{}
	}

	r.applyPodSpec(producer, existing, existing.Spec.Template)

	existing.Spec.Replicas = producer.Spec.Replicas
	existing.Spec.Selector = map[string]string{
		"app":              utils.MakeInstanceName(producer),
		"deploymentconfig": utils.DeploymentConfigName("prod", existing),
	}

	existing.Spec.Strategy.Type = appsv1.DeploymentStrategyTypeRecreate

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
	existing.Spec.Triggers[1].ImageChangeParams.ContainerNames = []string{"producer"}
	existing.Spec.Triggers[1].ImageChangeParams.From.Kind = "ImageStreamTag"
	existing.Spec.Triggers[1].ImageChangeParams.From.Name = utils.MakeInstanceName(producer) + "-parent:latest"

}

func (r *ReconcileProducer) configureMqtt(producer *simv1alpha1.SimulatorProducer, pod *corev1.PodTemplateSpec, messageType string) {

	pod.Spec.Containers[0].Command = []string{
		"java",
		"-Dvertx.cacheDirBase=/tmp",
		"-Dvertx.logger-delegate-factory-class-name=io.vertx.core.logging.SLF4JLogDelegateFactory",
		"-jar",
		"/build/simulator-mqtt/target/simulator-mqtt-app.jar"}

}

func (r *ReconcileProducer) configureHttp(producer *simv1alpha1.SimulatorProducer, pod *corev1.PodTemplateSpec, messageType string) {

	pod.Spec.Containers[0].Command = []string{
		"java",
		"-Dvertx.cacheDirBase=/tmp",
		"-Dvertx.logger-delegate-factory-class-name=io.vertx.core.logging.SLF4JLogDelegateFactory",
		"-jar",
		"/build/simulator-http/target/simulator-http-app.jar"}

}

func mqttVariables(simulatorName string) []corev1.EnvVar {

	return []v1.EnvVar{
		{Name: "HONO_MQTT_HOST", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "mqttAdapter.host"}}},
		{Name: "HONO_MQTT_PORT", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "mqttAdapter.port"}}},
	}

}

func httpVariables(simulatorName string) []corev1.EnvVar {

	return []v1.EnvVar{
		{Name: "HONO_HTTP_URL", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: simulatorName}, Key: "httpAdapter.url"}}},
		{Name: "DEVICE_PROVIDER", Value: "VERTX"},
		{Name: "VERTX_POOLED_BUFFERS", Value: "true"},
		{Name: "VERTX_RECREATE_CLIENT", Value: "120000"},
		{Name: "VERTX_KEEP_ALIVE", Value: "true"},
	}

}
