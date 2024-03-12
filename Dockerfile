FROM golang:1.21


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /blood-for-life-api

CMD ["/blood-for-life-api"]