FROM golang:1.14

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build main.go

CMD ["./main"]