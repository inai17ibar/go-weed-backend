# Dockerfile (backend/Dockerfile)

# ベースイメージ
FROM golang:1.21

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースをコピー
COPY . .

# 開放するポートを指定
EXPOSE 8081

# バイナリをビルド
RUN go build -o main ./cmd/go-weed

# 起動コマンド
CMD [ "./main" ]
