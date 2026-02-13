FROM golang:1.26.0-alpine as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /src/main .

CMD ["./main"]