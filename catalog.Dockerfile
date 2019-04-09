FROM quay.io/operator-framework/upstream-registry-builder

EXPOSE 50051
ENTRYPOINT ["/build/bin/registry-server"]
CMD ["--database", "bundles.db"]

COPY manifests manifests
RUN ./bin/initializer -o /build/bundles.db
