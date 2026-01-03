FROM golang:1.23 as builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/app

FROM alpine:3.20
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app
COPY --from=builder /bin/app /app/app
COPY configs/ /app/configs/
EXPOSE 8080
ENTRYPOINT ["/app/app"]


