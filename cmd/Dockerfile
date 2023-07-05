FROM golang:alpine

WORKDIR /go/src/postgis-vs-elasticsearch-performance

COPY ../app .

RUN go mod download

RUN go build -o myapp cmd/main.go

EXPOSE 8080

CMD ["./myapp"]