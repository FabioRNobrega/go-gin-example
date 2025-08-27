FROM golang:1.23

WORKDIR /app

# Ferramentas úteis + Air para hot reload
RUN apt-get update && apt-get install -y git vim curl && \
    go install github.com/cosmtrek/air@v1.52.0

CMD ["air"]
