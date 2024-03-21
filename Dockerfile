FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=linux go build -o /core-service cmd/core-service/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates

ENV APP_HOME=/home/bibliobreeze
ENV APP_USER=2100

RUN addgroup --system --gid ${APP_USER} bibliobreeze && \
    adduser --system --disabled-password --disabled-login \
       --home ${APP_HOME} --gecos "App User" --ingroup bibliobreeze \
       --uid ${APP_USER} bibliobreeze && \
    chown -R bibliobreeze:root ${APP_HOME} && \
    chmod -R 0775 ${APP_HOME}

COPY --from=build-stage --chown=${APP_USER} /core-service ./core-service
RUN chmod +x ./core-service

EXPOSE 8080

USER ${APP_USER}

ENTRYPOINT ["./core-service"]