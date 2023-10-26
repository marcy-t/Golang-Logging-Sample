FROM envoyproxy/envoy-alpine:v1.21.6 as local_envoy
COPY ./deployment/proxy/envoy.dev.yaml /etc/envoy/envoy.dev.yaml
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/usr/local/bin/envoy", "-c", "/etc/envoy/envoy.dev.yaml", "-l", "trace"]

