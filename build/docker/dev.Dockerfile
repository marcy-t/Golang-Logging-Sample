FROM golang:1.21.3-alpine3.18 as local

# Import environment variables from envfile.
# ARG GOPATH
# ARG STRIPE_SECRET_TEST
# ARG GCP_PROJECT_ID
# ARG SYSTEM_TIMEZONE
# ARG SYSTEM_LANG
# ARG DOCKER_USER
# ARG DOCKER_PASSWORD
# ARG DOCKER_GROUP
# ARG POSTGRES_USER
# ARG POSTGRES_PASSWORD
# ARG REDIS_USER
# ARG REDIS_PASSWORD
# ARG BEAST_HOST
# ARG BEAST_PORT

# # Set environment variables.
# ENV POSTGRES_USER $POSTGRES_USER
# ENV POSTGRES_PASSWORD $POSTGRES_PASSWORD
# ENV REDIS_USER $REDIS_USER
# ENV REDIS_PASSWORD $REDIS_PASSWORD
# ENV BEAST_HOST $BEAST_HOST
# ENV BEAST_PORT $BEAST_PORT
# ENV GOPATH $GOPATH
# ENV STRIPE_SECRET_TEST $STRIPE_SECRET_TEST
# ENV GCP_PROJECT_ID $GCP_PROJECT_ID
# ENV TZ $SYSTEM_TIMEZONE
# ENV LANG $SYSTEM_LANG
# ENV LINTER_VERSION=1.45.2

# Add apk for protobuf and its dependencies, and system packages.
RUN set -ex && \
    apk --no-cache update && \
    apk add git \
    binutils-gold \
    gcc \
    g++ \
    bash \
    tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.51.2 && \
    mv ./bin/golangci-lint /usr/bin/ && rm -r ./bin && \
    apk del tzdata

# Set permission to non-root user, and switch to non-root user.
WORKDIR /app

# Copy from source to the container.
COPY ./cmd/instagram_insight /app/cmd
COPY ./domain /app/domain
COPY ./infrastructure /app/infrastructure
COPY ./interfaces /app/interfaces
COPY ./usecase /app/usecase
COPY ./go.mod /app
COPY ./go.sum /app

COPY ./go.sum /app

# Install protofile compiler.
RUN go install github.com/cespare/reflex@latest && \
    go install github.com/golang/mock/mockgen@latest

# Copy from pkg source to the container.
COPY ./pkg /app/pkg

# Open the port for  server.
EXPOSE 8080


