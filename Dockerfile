# Build the project to binary file
FROM golang as builder

WORKDIR /app

COPY . /app

RUN go mod download
RUN go build -o main .

# Copy the binary file to the image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

CMD ["./main"]