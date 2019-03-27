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

package simulator

import (
	"context"
	"strconv"

	"github.com/ctron/operator-tools/pkg/install/core/secret"

	"github.com/ctron/operator-tools/pkg/install/core/configmap"

	"k8s.io/apimachinery/pkg/api/resource"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/ctron/iot-simulator-operator/pkg/install/prometheus"

	"github.com/ctron/operator-tools/pkg/install/openshift/route"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ctron/operator-tools/pkg/install/core/service"

	"github.com/ctron/operator-tools/pkg/install"

	"github.com/ctron/operator-tools/pkg/install/openshift/dc"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ctron/operator-tools/pkg/install/rbac/role"
	"github.com/ctron/operator-tools/pkg/install/rbac/rolebinding"

	buildv1 "github.com/openshift/api/build/v1"

	"github.com/ctron/operator-tools/pkg/install/openshift/build"

	"github.com/ctron/operator-tools/pkg/install/mixin"

	"github.com/ctron/operator-tools/pkg/install/core/serviceaccount"
	"github.com/ctron/operator-tools/pkg/install/openshift/imagestream"
	"github.com/ctron/operator-tools/pkg/recon"

	simv1alpha1 "github.com/ctron/iot-simulator-operator/pkg/apis/simulator/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_simulator")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSimulator{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {

	c, err := controller.New("simulator-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &simv1alpha1.Simulator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
		IsController: true, OwnerType: &simv1alpha1.Simulator{},
	})

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSimulator{}

type ReconcileSimulator struct {
	client client.Client
	scheme *runtime.Scheme
}

func instanceName(instance metav1.Object, basename string) string {
	return instance.GetName() + "-" + basename
}

func (r *ReconcileSimulator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Producer")

	// Fetch the Simulator instance
	instance := &simv1alpha1.Simulator{}
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

	// name := instance.Name

	rec := recon.NewContext(context.TODO(), request, r.client, r.scheme)

	ownerFn := mixin.ControllerOwner(instance)
	sharedOwnerFn := mixin.ObjectOwner(instance)

	// image streams

	rec.Process(imagestream.EmptyImageStream(instanceName(instance, "iot-simulator-base"), ownerFn))
	rec.Process(imagestream.EmptyImageStream(instanceName(instance, "iot-simulator-parent"), ownerFn))

	rec.Process(imagestream.EmptyImageStream("iot-simulator-console", sharedOwnerFn))

	rec.Process(imagestream.DockerImageStream("centos", "7", "docker.io/centos:7", sharedOwnerFn))
	rec.Process(imagestream.DockerImageStream("fedora", "29", "docker.io/fedora:29", sharedOwnerFn))

	// iot-simulator-console

	rec.Process(serviceaccount.ServiceAccount("iot-simulator-console", mixin.Mix(
		sharedOwnerFn,
		mixin.AppendAnnotationMixin("serviceaccounts.openshift.io/oauth-redirectreference.primary", `{"kind":"OAuthRedirectReference","apiVersion":"v1","reference":{"kind":"Route","name":"iot-simulator-console"}}`),
	)))
	rec.Process(role.WithRules("iot-simulator-console", []rbacv1.PolicyRule{
		{APIGroups: []string{""}, Resources: []string{"pods"}, Verbs: []string{"get", "list", "watch"}},
		{APIGroups: []string{"apps.openshift.io"}, Resources: []string{"deploymentconfigs"}, Verbs: []string{"get", "list", "watch"}},
	}, sharedOwnerFn))
	rec.Process(rolebinding.ForServiceAccount("iot-simulator-console", "iot-simulator-console", "iot-simulator-console", sharedOwnerFn))

	// prometheus

	rec.Process(serviceaccount.ServiceAccount("iot-simulator-prometheus", sharedOwnerFn))
	rec.Process(role.WithRules("iot-simulator-prometheus", []rbacv1.PolicyRule{
		{APIGroups: []string{""}, Resources: []string{"services", "endpoints", "pods"}, Verbs: []string{"get", "list", "watch"}},
		{APIGroups: []string{""}, Resources: []string{"configmaps"}, Verbs: []string{"get"}},
	}, sharedOwnerFn))
	rec.Process(rolebinding.ForServiceAccount("iot-simulator-prometheus", "iot-simulator-prometheus", "iot-simulator-prometheus", sharedOwnerFn))

	// build configs

	rec.Process(build.ReconcileBuildConfigSimple("iot-simulator-base", func(config *buildv1.BuildConfig) error {

		build.SetDockerStrategyFromImageStream(config, "centos:7")
		build.SetGitSource(config, "https://github.com/ctron/hono-simulator", "develop")
		config.Spec.Source.ContextDir = "containers/base"
		build.SetOutputImageStream(config, instanceName(instance, "iot-simulator-base")+":latest")
		build.EnableDefaultTriggers(config)

		return nil
	}, mixin.Mix(
		sharedOwnerFn,
	)))

	rec.Process(build.ReconcileBuildConfigSimple("iot-simulator-parent", func(config *buildv1.BuildConfig) error {

		build.SetDockerStrategyFromImageStream(config, instanceName(instance, "iot-simulator-base")+":latest")
		build.SetGitSource(config, "https://github.com/ctron/hono-simulator", "develop")
		build.SetOutputImageStream(config, instanceName(instance, "iot-simulator-parent")+":latest")
		build.EnableDefaultTriggers(config)

		return nil
	}, mixin.Mix(
		sharedOwnerFn,
	)))

	rec.Process(build.ReconcileBuildConfigSimple("iot-simulator-console", func(config *buildv1.BuildConfig) error {

		build.SetDockerStrategyFromImageStream(config, "fedora:29")
		build.SetGitSource(config, "https://github.com/ctron/iot-simulator-console", "develop")
		build.SetOutputImageStream(config, "iot-simulator-console:latest")
		build.EnableDefaultTriggers(config)

		config.Spec.Strategy.DockerStrategy.DockerfilePath = "Dockerfile.s2i"

		return nil
	}, mixin.Mix(
		sharedOwnerFn,
	)))

	// deployments

	rec.Process(dc.ReconcileDeploymentConfigSimple("iot-simulator-console", func(dc *appsv1.DeploymentConfig) error {

		if dc.Labels == nil {
			dc.Labels = make(map[string]string)
		}
		dc.Labels["app"] = "iot-simulator-console"
		dc.Labels["deploymentconfig"] = dc.Name

		dc.Spec.Selector = map[string]string{
			"app":              dc.Labels["app"],
			"deploymentconfig": dc.Labels["deploymentconfig"],
		}

		dc.Spec.Replicas = 1

		// template

		dc.Spec.Strategy.Type = appsv1.DeploymentStrategyTypeRolling

		if dc.Spec.Template == nil {
			dc.Spec.Template = &corev1.PodTemplateSpec{}
		}

		if dc.Spec.Template.ObjectMeta.Labels == nil {
			dc.Spec.Template.ObjectMeta.Labels = make(map[string]string)
		}

		dc.Spec.Template.ObjectMeta.Labels["app"] = dc.Labels["app"]
		dc.Spec.Template.ObjectMeta.Labels["deploymentconfig"] = dc.Labels["deploymentconfig"]

		// template spec

		dc.Spec.Template.Spec.ServiceAccountName = "iot-simulator-console"

		// containers

		if len(dc.Spec.Template.Spec.Containers) != 2 {
			dc.Spec.Template.Spec.Containers = make([]corev1.Container, 2)
		}

		// container - console

		dc.Spec.Template.Spec.Containers[0].Name = "console"
		dc.Spec.Template.Spec.Containers[0].Image = "iot-simulator-console"
		dc.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullAlways

		dc.Spec.Template.Spec.Containers[0].Env = []corev1.EnvVar{
			{Name: "GIN_MODE", Value: "release"},
			{Name: "PROMETHEUS_HOST", Value: "prometheus-operated"},
			install.EnvVarNamespace("NAMESPACE"),
		}

		dc.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
			{
				ContainerPort: 8080,
				Name:          "ui",
				Protocol:      corev1.ProtocolTCP,
			},
		}

		// container - oauth proxy

		dc.Spec.Template.Spec.Containers[1].Name = "oauth-proxy"
		dc.Spec.Template.Spec.Containers[1].Image = "openshift3/oauth-proxy"
		dc.Spec.Template.Spec.Containers[1].ImagePullPolicy = corev1.PullIfNotPresent

		dc.Spec.Template.Spec.Containers[1].Args = []string{
			"--https-address=:8443",
			"--provider=openshift",
			"--openshift-service-account=iot-simulator-console",
			"--upstream=http://localhost:8080",
			"--tls-cert=/etc/tls/private/tls.crt",
			"--tls-key=/etc/tls/private/tls.key",
			"--cookie-secret=SECRET",
		}

		dc.Spec.Template.Spec.Containers[1].Ports = []corev1.ContainerPort{
			{
				ContainerPort: 8443,
				Name:          "proxy",
				Protocol:      corev1.ProtocolTCP,
			},
		}

		if len(dc.Spec.Template.Spec.Containers[1].VolumeMounts) != 1 {
			dc.Spec.Template.Spec.Containers[1].VolumeMounts = make([]corev1.VolumeMount, 1)
		}
		dc.Spec.Template.Spec.Containers[1].VolumeMounts[0].Name = "proxy-tls"
		dc.Spec.Template.Spec.Containers[1].VolumeMounts[0].MountPath = "/etc/tls/private"

		// triggers

		if len(dc.Spec.Triggers) != 2 {
			dc.Spec.Triggers = make([]appsv1.DeploymentTriggerPolicy, 2)
		}

		dc.Spec.Triggers[0].Type = appsv1.DeploymentTriggerOnConfigChange
		dc.Spec.Triggers[1].Type = appsv1.DeploymentTriggerOnImageChange
		if dc.Spec.Triggers[1].ImageChangeParams == nil {
			dc.Spec.Triggers[1].ImageChangeParams = &appsv1.DeploymentTriggerImageChangeParams{}
		}
		dc.Spec.Triggers[1].ImageChangeParams.Automatic = true
		dc.Spec.Triggers[1].ImageChangeParams.ContainerNames = []string{dc.Spec.Template.Spec.Containers[0].Name}
		dc.Spec.Triggers[1].ImageChangeParams.From.Kind = "ImageStreamTag"
		dc.Spec.Triggers[1].ImageChangeParams.From.Name = "iot-simulator-console:latest"

		// volumes

		if len(dc.Spec.Template.Spec.Volumes) != 1 {
			dc.Spec.Template.Spec.Volumes = make([]corev1.Volume, 1)
		}
		dc.Spec.Template.Spec.Volumes[0].Name = "proxy-tls"
		if dc.Spec.Template.Spec.Volumes[0].Secret == nil {
			dc.Spec.Template.Spec.Volumes[0].Secret = &corev1.SecretVolumeSource{}
		}
		dc.Spec.Template.Spec.Volumes[0].Secret.SecretName = "iot-simulator-console-tls"

		// return

		return nil

	}, mixin.Mix(
		sharedOwnerFn,
	)))

	rec.Process(service.Service("iot-simulator-console", map[string]string{
		"app":              "iot-simulator-console",
		"deploymentconfig": "iot-simulator-console",
	}, func(service *corev1.Service) error {

		service.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "ui",
				Port:       8443,
				TargetPort: intstr.FromInt(8443),
				Protocol:   corev1.ProtocolTCP,
			},
		}

		return nil
	}, mixin.Mix(
		sharedOwnerFn,
		mixin.AppendAnnotationMixin("service.alpha.openshift.io/serving-cert-secret-name", "iot-simulator-console-tls"),
	)))

	rec.Process(route.ReencryptRoute("iot-simulator-console", "iot-simulator-console", intstr.FromString("ui"), sharedOwnerFn))

	// prometheus

	rec.Process(prometheus.ReconcilePrometheusSimple("iot-simulator-prometheus", func(prom *promv1.Prometheus) error {

		prom.Spec.ServiceAccountName = "iot-simulator-prometheus"
		prom.Spec.ServiceMonitorSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"metrics": "iot-simulator",
			},
		}

		prom.Spec.Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: *resource.NewQuantity(512*1024*1024 /* 512Mi */, resource.BinarySI),
			},
		}

		if prom.Spec.SecurityContext == nil {
			prom.Spec.SecurityContext = &corev1.PodSecurityContext{}
		}

		return nil
	}, sharedOwnerFn))

	rec.Process(prometheus.ReconcileServiceMonitorSimple("iot-simulator-prometheus", func(monitor *promv1.ServiceMonitor) error {

		monitor.Spec.Selector.MatchLabels = map[string]string{
			"metrics": "iot-simulator",
		}

		if len(monitor.Spec.Endpoints) != 1 {
			monitor.Spec.Endpoints = make([]promv1.Endpoint, 1)
		}

		monitor.Spec.Endpoints[0].Port = "metrics"
		monitor.Spec.Endpoints[0].Path = "/metrics"
		monitor.Spec.Endpoints[0].Interval = "1s"

		return nil
	}, mixin.Mix(
		sharedOwnerFn,
		mixin.AppendLabelMixin("metrics", "iot-simulator"),
	)))

	// endpoint information

	rec.Process(configmap.ReconcileConfigMapSimple(instance.Name, func(configMap *corev1.ConfigMap) error {
		if configMap.Data == nil {
			configMap.Data = make(map[string]string)
		}

		configMap.Data["endpoint.host"] = instance.Spec.Endpoint.Messaging.Host
		configMap.Data["endpoint.port"] = strconv.Itoa(instance.Spec.Endpoint.Messaging.Port)

		configMap.Data["deviceRegistry.url"] = instance.Spec.Endpoint.Registry.URL

		configMap.Data["mqttAdapter.host"] = instance.Spec.Endpoint.Adapters.MQTT.Host
		configMap.Data["mqttAdapter.port"] = strconv.Itoa(instance.Spec.Endpoint.Adapters.MQTT.Port)

		configMap.Data["httpAdapter.url"] = instance.Spec.Endpoint.Adapters.HTTP.URL

		return nil
	}, ownerFn))

	rec.Process(secret.ReconcileSecretSimple(instance.Name, func(secret *corev1.Secret) error {
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}

		secret.Data["endpoint.username"] = []byte(instance.Spec.Endpoint.Messaging.User)
		secret.Data["endpoint.password"] = []byte(instance.Spec.Endpoint.Messaging.Password)

		if len(instance.Spec.Endpoint.Messaging.CACertificate) > 0 {
			secret.Data["messaging.ca.crt"] = instance.Spec.Endpoint.Messaging.CACertificate
		}

		return nil
	}, ownerFn))

	return rec.Result()
}
