FROM golang:latest as build

WORKDIR /go/src/github.com/yurakawa/sbi-ipo-cp-saver
ADD . /go/src/github.com/yurakawa/sbi-ipo-cp-saver

RUN go mod download
RUN CGO_ENABLED=0 go build -v -o main .

FROM chromedp/headless-shell:latest
COPY --from=build /go/src/github.com/yurakawa/sbi-ipo-cp-saver/main /
CMD ["/main"]