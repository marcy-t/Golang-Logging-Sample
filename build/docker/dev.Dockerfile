FROM golang:1.21.3-alpine3.18 as local

# Import environment variables from envfile.
ARG MYSQL_HOST
ENV MYSQL_HOST $MYSQL_HOST
ARG MYSQL_DATABASE
ENV MYSQL_DATABASE $MYSQL_DATABASE
ARG MYSQL_USER
ENV MYSQL_USER $MYSQL_USER
ARG MYSQL_PASSWORD
ENV MYSQL_PASSWORD $MYSQL_PASSWORD
ARG TZ
ENV TZ $TZ

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

# Install dependencies for this project.
RUN go mod download

# Open the port for  server.
EXPOSE 8080


