FROM golang:1.23

WORKDIR /app

# Instala ferramentas e o 'air' para hot reload
RUN apt-get update && apt-get install -y git vim curl && \
    go install github.com/cosmtrek/air@v1.52.0

CMD ["air"]
