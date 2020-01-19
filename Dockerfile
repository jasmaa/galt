FROM golang:1.13-alpine

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build cmd/app.go

CMD [ "./app" ]