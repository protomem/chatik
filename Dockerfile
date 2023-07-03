FROM golang:alpine AS builder

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./build/ ./cmd/chatik 


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/build/ ./build/

CMD [ "./build/chatik" ]
