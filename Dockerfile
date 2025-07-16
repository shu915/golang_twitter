# ベースイメージ（開発用ツール込み）
FROM golang:1.24.5-alpine3.22 AS development

WORKDIR /app

# Airをインストール
RUN go install github.com/air-verse/air@latest

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードはボリュームマウントで同期するためコピーしない

EXPOSE 8080

# Airを実行
CMD ["air", "-c", ".air.toml"]


# ビルド用ステージ
FROM development AS builder

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN go build -o main /app/main.go


# 本番用ステージ（最小イメージ）
FROM alpine:3.22 AS production

WORKDIR /app

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates

# バイナリをコピー
COPY --from=builder /app/main .

EXPOSE 8080

# バイナリを実行
CMD ["./main"]