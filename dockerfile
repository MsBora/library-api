#Stage 1
FROM golang:1.24-alpine3.22 AS builder

RUN apk add --no-cache git


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o /library-api .

#Stage 2
FROM alpine:latest

COPY --from=builder /library-api /library-api

EXPOSE 8080

CMD ["/library-api"]

