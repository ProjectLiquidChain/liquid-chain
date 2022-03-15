FROM golang:1.15 AS builder

WORKDIR /liquid-chain
ADD go.mod go.sum /liquid-chain/
RUN go mod download
ADD . /liquid-chain
RUN cd /liquid-chain/cmd && \
  go build -a -installsuffix nocgo -tags builtin_static -o /lqc-node .

ENTRYPOINT ["/lqc-node"]
