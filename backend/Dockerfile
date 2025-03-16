FROM golang:1.24.1-bookworm

WORKDIR /app

COPY go.mod .
COPY go.sum .
#RUN sudo apt install gcc
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o main .

EXPOSE 8080

CMD ["./main"]