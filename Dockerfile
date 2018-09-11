FROM golang:1.11


WORKDIR /go/src/github.com/jaholmes/kubclient
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["kubclient"]