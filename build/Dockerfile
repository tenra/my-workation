FROM golang:1.16

COPY go.mod /my-workation/go.mod
COPY go.sum /my-workation/go.sum
WORKDIR /my-workation
RUN ls
RUN go mod download
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jinzhu/gorm
RUN go get github.com/joho/godotenv

COPY . /my-workation

CMD go run main.go
