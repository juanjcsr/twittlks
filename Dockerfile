FROM golang:1.17-alpine as builder

WORKDIR /app/

COPY go.mod go.sum /app/
RUN go mod download

COPY . .

RUN  CGO_ENABLED=0 go build -o twitlks

FROM alpine
COPY --from=builder /app/twitlks /
ENTRYPOINT ["/twitlks"]
