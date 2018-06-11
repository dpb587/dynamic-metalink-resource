FROM golang:1.10 as resource
WORKDIR /go/src/github.com/dpb587/dynamic-metalink-resource
COPY . .
ENV CGO_ENABLED=0
RUN mkdir -p /opt/resource
RUN git rev-parse HEAD | tee /opt/resource/version
RUN go build -o /opt/resource/check check/*.go
RUN go build -o /opt/resource/in in/*.go

FROM alpine:3.4
RUN apk --no-cache add bash ca-certificates coreutils curl git gnupg openssh-client jq
COPY --from=resource /opt/resource /opt/resource
