FROM golang:1.24.1-alpine3.21

WORKDIR /go/src/neo-inu

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install github.com/air-verse/air@latest

CMD [ "go", "tool", "air" ]
