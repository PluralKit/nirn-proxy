FROM --platform=$BUILDPLATFORM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags '-extldflags "-static"' -tags timetzdata -buildvcs=false

FROM alpine:latest
COPY --from=app-builder /go/src/app/nirn-proxy /nirn-proxy
EXPOSE 9000
EXPOSE 8080
ENTRYPOINT ["/nirn-proxy"]
