# Stage 1: build
FROM golang:1.23 AS builder
WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./main.go

# Stage 2: runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# Copy only the binary
COPY --from=builder /app/app .

# The app expects RSA_PRIVATE_KEY as env injected by Kubernetes
ENV RSA_PRIVATE_KEY=""

EXPOSE 8082
CMD ["./app"]
