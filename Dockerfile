# Build the project to binary file
FROM golang as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Copy the binary file to the image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config.yml .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/logs ./logs

CMD ["./main"]