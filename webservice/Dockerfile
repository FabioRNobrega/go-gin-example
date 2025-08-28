# Stage 1: build
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./main.go

# Stage 2: minimal run
FROM scratch
WORKDIR /app

# Copy binary
COPY --from=builder /app/app .

# Copy required folders for frontend
COPY static ./static
COPY templates ./templates

EXPOSE 8080
CMD ["./app"]
