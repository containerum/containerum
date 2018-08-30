FROM golang:1.10-alpine as builder
ENV PATH=$PATH:$GOPATH/bin
WORKDIR $GOPATH/src/github.com/containerum/containerum/embark
COPY . .
RUN go install -v ./cmd/embark/

# ---
FROM alpine:3.7
COPY --from=builder /go/bin/embark /embark
COPY --from=builder /go/src/github.com/containerum/containerum/embark/charts/containerum  /charts/containerum

ENTRYPOINT ["/embark", "--kube-config", "/etc/kube/config", "--debug", "--namespace", "containerum", "--dir", "/charts/containerum"]
