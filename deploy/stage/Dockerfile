FROM golang:alpine AS builder

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./build/ ./cmd/chatik 


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/ ./build/
COPY --from=builder /app/configs/ ./configs/

CMD [ "./build/chatik", "-conf", "./configs/stage/app.yaml" ]
