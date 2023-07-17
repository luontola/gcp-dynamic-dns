FROM golang:1.20 AS builder

# download dependencies
WORKDIR /go/src/app
COPY src/app/go.mod src/app/go.sum /go/src/app/
RUN go mod download && go mod verify

ENV GOFLAGS -tags=netgo

# build the app
COPY src/app /go/src/app
RUN go test -v ./...
RUN go install -v -ldflags "-linkmode external -extldflags -static" .

# ------------------------------------------------------------

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /
USER 1000
ENTRYPOINT ["/app"]

COPY --from=builder /go/bin/app /app
