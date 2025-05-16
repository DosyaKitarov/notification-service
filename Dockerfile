FROM golang:1.23

WORKDIR /app

RUN apt-get update && apt-get install -y postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/notification/"]