FROM golang

WORKDIR /app

COPY . /app

RUN go mod download
RUN go build -o main .

EXPOSE $PORT

CMD ["./main"]