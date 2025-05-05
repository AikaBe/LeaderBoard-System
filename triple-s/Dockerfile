FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o triple-s main.go

EXPOSE 9000

CMD ["./triple-s","--dir=/data"]