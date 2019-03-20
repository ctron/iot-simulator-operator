package common

import (
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
)

func ApplyProbe(probe *corev1.Probe) *corev1.Probe {
	if probe == nil {
		probe = &corev1.Probe{}
	}
	probe.Exec = nil
	probe.TCPSocket = nil
	probe.HTTPGet = &corev1.HTTPGetAction{
		Path:   "/health",
		Port:   intstr.FromInt(8081),
		Scheme: corev1.URISchemeHTTP,
	}
	return probe
}
