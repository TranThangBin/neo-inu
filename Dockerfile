FROM golang:1.23.4-alpine3.20

WORKDIR /go/src/neo-inu

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install github.com/air-verse/air@latest

CMD $(go env GOPATH)/bin/air
