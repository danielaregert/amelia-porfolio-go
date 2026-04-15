# Runtime image para el binario ya compilado para linux/arm64.
# El binario se compila en el host con cmd/build-pi.sh.
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

ENV TZ=America/Argentina/Buenos_Aires
WORKDIR /app

COPY dist/porfolio-amelia /app/porfolio-amelia
RUN chmod +x /app/porfolio-amelia

EXPOSE 8090

# pb_data va por volumen (definido en docker-compose.yml).
CMD ["/app/porfolio-amelia", "serve", "--http=0.0.0.0:8090"]
