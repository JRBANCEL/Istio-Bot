FROM golang:1.12 as build

WORKDIR /go/src/github.com/JRBANCEL/Istio-Bot
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM gcr.io/distroless/base
COPY --from=build /go/bin/Istio-Bot /
CMD ["/Istio-Bot"]
