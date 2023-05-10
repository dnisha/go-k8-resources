FROM golang:1.20-alpine

WORKDIR /app

COPY . .

EXPOSE 8080

RUN go mod download
RUN go build -o main main.go

EXPOSE 8080

CMD ["/app/main"]
