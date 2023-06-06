# syntax=docker/dockerfile:experimental
# ---
FROM golang:1.20 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV NAMESPACES_TO_IGNORE=kube-system
ENV VULNERABILITIES_TO_IGNORE=
ENV KUBECLARITY_API_URL=http://kubeclarity-kubeclarity.kubeclarity.svc:8080
ENV KEEL_WEBHOOK_URL=https://keel-webhook.trittale.svc

WORKDIR /work
COPY api /work/api
COPY di /work/di
COPY dto /work/dto
COPY service /work/service
COPY go.mod /work/go.mod
COPY go.sum /work/go.sum
COPY LICENSE /work/LICENSE
COPY main.go /work/main.go
COPY README.md /work/README.md


RUN --mount=type=cache,target=/root/.cache/go-build,sharing=private \
  go build -o bin/scanyourkube-management-job .

FROM alpine:latest AS create-utils
RUN adduser -u 10001 scratchuser -D

# ---
FROM scratch AS run
COPY --from=create-utils etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /work/bin/scanyourkube-management-job /usr/local/bin/
COPY --from=create-utils /etc/passwd /etc/passwd
USER scratchuser
CMD ["scanyourkube-management-job"]
