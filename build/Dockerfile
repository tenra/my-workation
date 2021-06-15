FROM golang:1.16

COPY go.mod /my-workation/go.mod
COPY go.sum /my-workation/go.sum
WORKDIR /my-workation
RUN ls
RUN go mod download

COPY . /my-workation

CMD go run main.go
