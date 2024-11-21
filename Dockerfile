FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

COPY go.mod go.sum ./
RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o /app/main .

FROM alpine:latest

LABEL author="Balaji Shettigar"

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]

CMD ["start"]