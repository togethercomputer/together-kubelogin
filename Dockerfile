FROM --platform=$BUILDPLATFORM golang:1.26.10@sha256:672c944cfd07599b5f1bb179869270cdc94c6254e41db4d431259693cb286b82 AS builder

WORKDIR /builder

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
COPY main.go .
COPY pkg pkg

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o kubectl-together_login

FROM gcr.io/distroless/base-debian13
COPY --from=builder /builder/kubectl-together_login /
ENTRYPOINT ["/kubectl-together_login"]
