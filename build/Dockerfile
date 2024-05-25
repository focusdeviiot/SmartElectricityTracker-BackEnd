FROM golang:1.18-alpine

WORKDIR /app

COPY ../go.mod .
COPY ../go.sum .

RUN go mod download

COPY . .

COPY ../configs ./configs

RUN go build -o main ../cmd/smart_electricity_tracker_backend

CMD ["./main"]