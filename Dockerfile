FROM golang:1.25 AS builder
WORKDIR /src

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./main.go ./
COPY ./help_handler.go ./
COPY ./hello_handler.go ./
COPY ./message_router.go  ./

# プラットフォームに応じてGOOSとGOARCHを設定
RUN case ${TARGETPLATFORM} in \
        "linux/amd64")  export GOOS=linux GOARCH=amd64 ;; \
        "linux/arm64")  export GOOS=linux GOARCH=arm64 ;; \
        "linux/arm/v7") export GOOS=linux GOARCH=arm GOARM=7 ;; \
        *) echo "Unsupported platform: ${TARGETPLATFORM}" && exit 1 ;; \
    esac && \
    CGO_ENABLED=0 go build -o /tacokumo-bot ./tacokumo-bot

FROM scratch
CMD ["/tacokumo-bot"]
