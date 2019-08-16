# IoT Simulator Operator

This is the Kubernetes operator repository for the IoT simulator.

See https://github.com/ctron/hono-simulator for more information.

## Creating a simulator

First you need to create a simulator base, assuming that your cluster
apps base domain is `apps.your.cluster` and the project you
deployed EnMasse to the project `enmasse-infra` and the IoT simulator
to the project `iot-simulator`:

~~~yaml
kind: Simulator
apiVersion: iot.dentrassi.de/v1alpha1
metadata:
  name: hono1
spec:
  endpoint:
    adapters:
      http:
        url: https://iot-http-adapter-enmasse-infra.apps.your.cluster
      mqtt:
        host: iot-mqtt-adapter-enmasse-infra.apps.your.cluster
        port: 443
    messaging:
      caCertificate: <base64 encoded PKCS#1/PEM cert>
      host: messaging-<infraUUID>.enmasse-infra.svc
      port: 5671
      user: consumer
      password: foobar
    registry:
      url: https://device-registry.apps.your.cluster
~~~

Next you need a consumer, created for the IoT tenant `iot-simulator.iot`:

~~~yaml
kind: SimulatorConsumer
apiVersion: iot.dentrassi.de/v1alpha1
metadata:
  name: consumer1
spec:
  replicas: 1
  simulator: hono1
  tenant: iot-simulator.iot
  type: telemetry
~~~

Then you can create a producer:

~~~yaml
kind: SimulatorProducer
apiVersion: iot.dentrassi.de/v1alpha1
metadata:
  name: producer1
spec:
  numberOfDevices: 10
  protocol: http
  replicas: 1
  simulator: hono1
  tenant: iot-simulator.iot
  type: telemetry
~~~

## OpenShift

When running in OpenShift, the operator will automatically set up ImageStreams,
Builds and Routes.

### Routes

You should be able to see statistics on the Web UI
<https://iot-simulator-console-iot-simulator.apps.your.cluster/>, as soon as
you created producers and consumers.

### Builds

By default it will build the matching version from the original repository.

However you can use the `Simulator` custom resource, to let the operator
tweak the build. You may use this to try you own variations of the simulator.

~~~yaml
kind: Simulator
apiVersion: iot.dentrassi.de/v1alpha1
metadata:
  name: hono1
spec:
  simulator:
    builds:
      hono-simulator:
        git:
          # Full configuration
          uri: https://github.com/ctron/hono-simulator
          ref: develop
      iot-simulator-console:
        git:
          # Only change branch
          ref: develop
~~~

## Building for OLM

~~~
docker build -t docker.io/ctron/iot-simulator-source:latest -f catalog.Dockerfile .
~~~

Load with `oc apply -f`:

~~~yaml
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: iot-simulator-source
  namespace: default
spec:
  sourceType: grpc
  image: docker.io/ctron/iot-simulator-source:latest
~~~

Also see:

* https://github.com/operator-framework/operator-registry#building-a-catalog-of-operators-using-operator-registry
