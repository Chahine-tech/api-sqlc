FROM golang:1.22.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /api

EXPOSE 8080

CMD ["/api"]
