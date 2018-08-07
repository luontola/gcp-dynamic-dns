FROM golang:1.10

# tools
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# download dependencies
WORKDIR /go/src/app
COPY src/main/Gopkg.toml src/main/Gopkg.lock /go/src/app/
RUN dep ensure -v -vendor-only

# build the app
COPY src/main /go/src/app
RUN go install -v .

CMD ["app"]
