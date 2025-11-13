FROM golang:1.25 AS builder
WORKDIR /src

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./main.go ./
COPY ./bot.go ./
COPY ./help_handler.go ./
COPY ./hello_handler.go ./
COPY ./message_router.go  ./
RUN CGO_ENABLED=0 go build -o /tacokumo-bot .

FROM scratch
COPY --from=builder /tacokumo-bot /tacokumo-bot
# CA証明書をコピー
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /server /server
CMD ["/tacokumo-bot"]
