FROM envoyproxy/envoy-alpine:v1.21.6 as local_mysql_envoy
COPY ./deployment/proxy/mysql_envoy.dev.yaml /etc/envoy/mysql_envoy.dev.yaml
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/usr/local/bin/envoy", "-c", "/etc/envoy/mysql_envoy.dev.yaml", "-l", "trace"]
