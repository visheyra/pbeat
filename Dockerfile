FROM golang:1.10 as BUILD

RUN mkdir -p ${GOPATH}/src/github.com/elastic && \
  git clone https://github.com/elastic/beats ${GOPATH}/src/github.com/elastic/beats

RUN go get -u -v github.com/prometheus/prometheus/... && \
  go get -u -v github.com/golang/snappy/...

RUN go get -u -v github.com/golang/protobuf/proto/...

WORKDIR ${GOPATH}/src/github.com/visheyra/pbeat

COPY . .

RUN go install -x -v github.com/visheyra/pbeat

FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=BUILD /go/bin/pbeat /app

CMD ["/app/pbeat"]
