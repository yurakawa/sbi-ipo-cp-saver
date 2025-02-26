FROM golang:1.23.0-alpine as build

WORKDIR /go/src/github.com/yurakawa/sbi-ipo-cp-saver
ADD . /go/src/github.com/yurakawa/sbi-ipo-cp-saver

# キャッシュを活用するために、依存関係のダウンロードを分ける
COPY go.mod go.sum ./
RUN go mod download

# ソースコードを追加
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main .

FROM chromedp/headless-shell:latest
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && update-ca-certificates \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
COPY --from=build /go/src/github.com/yurakawa/sbi-ipo-cp-saver/main /

CMD ["/main"]
