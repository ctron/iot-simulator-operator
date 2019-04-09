# IoT Simulator Operator

This is the kubernetes operator repository for the IoT simulator.

See https://github.com/ctron/hono-simulator for more information.

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
