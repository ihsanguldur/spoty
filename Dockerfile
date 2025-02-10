FROM golang:1.23.4-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o spoty cmd/web/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/spoty .

COPY ./assets ./assets

COPY ./internal/app/templates ./internal/app/templates

RUN chmode +x /app/spoty

ENTRYPOINT [ "/app/spoty" ]