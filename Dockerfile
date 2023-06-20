FROM golang:1.20.2-alpine AS builder

WORKDIR /go/src/github.com/yurakawa/sbi-ipo-cp-saver
ADD . /go/src/github.com/yurakawa/sbi-ipo-cp-saver

RUN go mod download
RUN go build -v -o main


FROM alpine:latest

COPY --from=builder /go/src/github.com/yurakawa/sbi-ipo-cp-saver/main /bin/main

CMD ["/bin/main", "--help"]