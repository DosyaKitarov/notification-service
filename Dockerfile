FROM golang:1.23

WORKDIR /app

RUN apt-get update && apt-get install

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/notification/"]