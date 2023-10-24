# FROM envoyproxy/envoy-alpine:v1.28.0 as local_envoy
FROM envoyproxy/envoy-alpine:v1.14.1 as local_envoy
# COPY ./deployment/proxy/envoy.yaml /etc/envoy/envoy.yaml
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/usr/local/bin/envoy", "-c", "/etc/envoy/envoy.yaml", "-l", "trace"]

