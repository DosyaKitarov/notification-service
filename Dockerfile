FROM golang:1.23

WORKDIR /app

RUN apt-get update && apt-get install

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 8080

CMD ["/wait-for-it.sh", "kafka:9092", "--", "go", "run", "./cmd/notification/"]