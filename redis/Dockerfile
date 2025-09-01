FROM redis:7.2-alpine

# Copy config
COPY redis.conf /usr/local/etc/redis/

# Expose the port
ARG port
EXPOSE $port

CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]