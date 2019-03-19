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

var log = logf.Log.WithName("controller_producer")

var TRUE = true

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

	err = c.Watch(&source.Kind{Type: &appsv1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
		IsController: true, OwnerType: &simv1alpha1.SimulatorProducer{},
	})
	if err != nil {
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

func (r *ReconcileProducer) configureService(producer *simv1alpha1.SimulatorProducer, existing *v1.Service) {

	if existing.ObjectMeta.Labels == nil {
		existing.ObjectMeta.Labels = map[string]string{}
	}

	existing.ObjectMeta.Labels["app"] = utils.MakeHelmInstanceName(producer)
	existing.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("prod", existing)
	existing.ObjectMeta.Labels["metrics"] = utils.MakeHelmInstanceName(producer)
	existing.ObjectMeta.Labels["iot.simulator"] = producer.Spec.Simulator

	existing.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "metrics",
			Port:       8081,
			TargetPort: intstr.FromInt(8081),
		},
	}
	existing.Spec.Selector = map[string]string{
		"app":              utils.MakeHelmInstanceName(producer),
		"deploymentconfig": existing.Name,
	}

}

func (r *ReconcileProducer) configureDeploymentConfig(producer *simv1alpha1.SimulatorProducer, existing *appsv1.DeploymentConfig) {

	if existing.ObjectMeta.Labels == nil {
		existing.ObjectMeta.Labels = map[string]string{}
	}

	endpointConfigName := producer.Spec.EndpointSettings
	messageType := producer.Spec.Type
	if messageType == "" {
		messageType = "telemetry"
	}
	protocol := producer.Spec.Protocol
	if protocol == "" {
		protocol = simv1alpha1.ProtocolHttp
	}

	existing.ObjectMeta.Labels["app"] = utils.MakeHelmInstanceName(producer)
	existing.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("prod", existing)
	existing.ObjectMeta.Labels["iot.simulator.tenant"] = producer.Spec.Tenant
	existing.ObjectMeta.Labels["iot.simulator"] = producer.Spec.Simulator
	existing.ObjectMeta.Labels["iot.simulator.app"] = "producer"
	existing.ObjectMeta.Labels["iot.simulator.message.type"] = messageType
	existing.ObjectMeta.Labels["iot.simulator.producer.protocol"] = string(protocol)

	existing.Spec.Replicas = producer.Spec.Replicas
	existing.Spec.Selector = map[string]string{
		"app":              utils.MakeHelmInstanceName(producer),
		"deploymentconfig": utils.DeploymentConfigName("prod", existing),
	}

	existing.Spec.Strategy.Type = appsv1.DeploymentStrategyTypeRecreate

	if existing.Spec.Template == nil {
		existing.Spec.Template = &v1.PodTemplateSpec{}
	}

	if existing.Spec.Template.ObjectMeta.Labels == nil {
		existing.Spec.Template.ObjectMeta.Labels = make(map[string]string)
	}

	existing.Spec.Template.ObjectMeta.Labels["app"] = utils.MakeHelmInstanceName(producer)
	existing.Spec.Template.ObjectMeta.Labels["deploymentconfig"] = utils.DeploymentConfigName("prod", existing)
	existing.Spec.Template.ObjectMeta.Labels["iot.simulator.tenant"] = producer.Spec.Tenant
	existing.Spec.Template.ObjectMeta.Labels["iot.simulator"] = producer.Spec.Simulator
	existing.Spec.Template.ObjectMeta.Labels["iot.simulator.settings"] = producer.Spec.EndpointSettings

	// containers

	if len(existing.Spec.Template.Spec.Containers) != 1 {
		existing.Spec.Template.Spec.Containers = make([]corev1.Container, 1)
	}

	existing.Spec.Template.Spec.Containers[0].Name = "producer"
	existing.Spec.Template.Spec.Containers[0].Env = []v1.EnvVar{

		{Name: "MESSAGE_TYPE", Value: messageType},

		{Name: "PERIOD_MS", Value: "1000"},
		{Name: "NUM_DEVICES", Value: strconv.FormatUint(uint64(producer.Spec.NumberOfDevices), 10)},

		{Name: "HONO_TENANT", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.labels['iot.simulator.tenant']"}}},

		{Name: "DEVICE_REGISTRY_URL", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: endpointConfigName}, Key: "deviceRegistry.url"}}},
		{Name: "TLS_INSECURE", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: endpointConfigName}, Key: "tlsInsecure", Optional: &TRUE}}},
	}

	if producer.Spec.NumberOfThreads != nil {
		existing.Spec.Template.Spec.Containers[0].Env = append(existing.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{
			Name: "NUM_THREADS", Value: strconv.FormatUint(uint64(*producer.Spec.NumberOfThreads), 10),
		})
	}

	existing.Spec.Template.Spec.Containers[0].Ports = []v1.ContainerPort{
		{
			ContainerPort: 8081,
			Name:          "metrics",
			Protocol:      corev1.ProtocolTCP,
		},
	}

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
	existing.Spec.Triggers[1].ImageChangeParams.From.Name = utils.MakeHelmInstanceName(producer) + "-parent:latest"

	// now apply http specifics

	switch protocol {
	case simv1alpha1.ProtocolMqtt:
		existing.Spec.Template.Spec.Containers[0].Env = append(existing.Spec.Template.Spec.Containers[0].Env, mqttVariables(endpointConfigName)...)
		r.configureMqtt(producer, existing, messageType)
	default:
		existing.Spec.Template.Spec.Containers[0].Env = append(existing.Spec.Template.Spec.Containers[0].Env, httpVariables(endpointConfigName)...)
		r.configureHttp(producer, existing, messageType)
	}

}

func (r *ReconcileProducer) configureMqtt(producer *simv1alpha1.SimulatorProducer, existing *appsv1.DeploymentConfig, messageType string) {

	existing.Spec.Template.Spec.Containers[0].Command = []string{"java",
		"-Xmx1024m",
		"-Dvertx.cacheDirBase=/tmp",
		"-Dvertx.logger-delegate-factory-class-name=io.vertx.core.logging.SLF4JLogDelegateFactory",
		"-jar",
		"/build/simulator-mqtt/target/simulator-mqtt-app.jar"}

}

func (r *ReconcileProducer) configureHttp(producer *simv1alpha1.SimulatorProducer, existing *appsv1.DeploymentConfig, messageType string) {

	existing.Spec.Template.Spec.Containers[0].Command = []string{"java",
		"-Xmx1024m",
		"-Dvertx.cacheDirBase=/tmp",
		"-Dvertx.logger-delegate-factory-class-name=io.vertx.core.logging.SLF4JLogDelegateFactory",
		"-jar",
		"/build/simulator-http/target/simulator-http-app.jar"}

}

func mqttVariables(sec string) []corev1.EnvVar {

	return []v1.EnvVar{
		{Name: "HONO_MQTT_HOST", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: sec}, Key: "mqttAdapter.host"}}},
		{Name: "HONO_MQTT_PORT", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: sec}, Key: "mqttAdapter.port"}}},
	}
}

func httpVariables(sec string) []corev1.EnvVar {

	return []v1.EnvVar{
		{Name: "HONO_HTTP_URL", ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: sec}, Key: "httpAdapter.url"}}},
		{Name: "DEVICE_PROVIDER", Value: "VERTX"},
		{Name: "VERTX_POOLED_BUFFERS", Value: "true"},
		{Name: "VERTX_RECREATE_CLIENT", Value: "120000"},
		{Name: "VERTX_KEEP_ALIVE", Value: "true"},
	}
}
