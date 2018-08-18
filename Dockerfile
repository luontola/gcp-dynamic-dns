FROM golang:1.10 AS builder

# tools
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# download dependencies
WORKDIR /go/src/app
COPY src/app/Gopkg.toml src/app/Gopkg.lock /go/src/app/
RUN dep ensure -v -vendor-only
RUN go build -v all

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
