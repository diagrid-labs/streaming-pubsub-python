FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o publisher main.go

FROM alpine:3.21

WORKDIR /app

# Copy binary from the builder stage
COPY --from=builder /app/publisher .

ENTRYPOINT ["./publisher"]