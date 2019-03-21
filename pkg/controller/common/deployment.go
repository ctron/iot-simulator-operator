package common

import (
	"encoding/json"
	"os"

	"github.com/ctron/iot-simulator-operator/pkg/apis/simulator/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
)

var Endpoint v1alpha1.SimulatorEndpoint
var WatchSimulatorName string

func init() {
	WatchSimulatorName = os.Getenv("WATCH_SIMULATOR_NAME")

	Endpoint = v1alpha1.SimulatorEndpoint{}
	if err := json.Unmarshal([]byte(os.Getenv("SIMULATOR_CONFIG")), &Endpoint); err != nil {
		panic(err)
	}
}

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
