FROM golang:1.23.0-alpine as build

WORKDIR /go/src/github.com/yurakawa/sbi-ipo-cp-saver
ADD . /go/src/github.com/yurakawa/sbi-ipo-cp-saver

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main .

FROM chromedp/headless-shell:129.0.6668.12
COPY --from=build /go/src/github.com/yurakawa/sbi-ipo-cp-saver/main /
CMD ["/main"]
