FROM golang:1.23
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /goapp cmd/architect_shortener/main.go
CMD ["/goapp"]