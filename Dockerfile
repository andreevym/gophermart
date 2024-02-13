# Build stage
FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gofermart cmd/gofermart/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/gofermart .

CMD ["./gofermart"]
