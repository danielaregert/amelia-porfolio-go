# Deploy a Raspberry Pi 5

Ruta sugerida en el Pi: `~/services/amelia-porfolio/`

## Pre-requisitos

- Caddy corriendo en Docker, conectado a la red `caddy_net`.
- Cloudflare Tunnel apuntando `ameporfolio.danielaregert.com.ar` → `caddy:80`.
- Si la red de Caddy se llama distinto, editá `docker-compose.yml` (cambiá `caddy_net`).

## Pasos (corren en el Mac, salvo el último)

```bash
# 1. Compilar binario para el Pi (linux/arm64)
./cmd/build-pi.sh

# 2. Crear carpeta destino en el Pi
ssh danielapi@danielapi.local 'mkdir -p ~/services/amelia-porfolio/dist'

# 3. Copiar el binario
scp dist/porfolio-amelia danielapi@danielapi.local:~/services/amelia-porfolio/dist/

# 4. Copiar Dockerfile + compose
scp Dockerfile docker-compose.yml danielapi@danielapi.local:~/services/amelia-porfolio/

# 5. Copiar pb_data (DB sqlite + storage con todas las imágenes/video/PDF)
rsync -avh --delete pb_data/ danielapi@danielapi.local:~/services/amelia-porfolio/pb_data/

# 6. Agregar bloque al Caddyfile (en el Pi) y reload
#    Pegar deploy/Caddyfile.snippet al final del Caddyfile, luego:
#      docker exec caddy caddy reload --config /etc/caddy/Caddyfile

# 7. Levantar el contenedor (en el Pi)
ssh danielapi@danielapi.local
cd ~/services/amelia-porfolio
docker compose up -d --build
docker compose logs -f amelia-porfolio
```

## Verificación

```bash
# Desde el Pi (interno)
curl -I http://localhost:8090/        # debería responder via el contenedor
docker exec caddy curl -I http://amelia-porfolio:8090/   # caddy → app

# Desde fuera
curl -I https://ameporfolio.danielaregert.com.ar/
```

## Backups

`pb_data/` es el único estado importante. Para snapshot:

```bash
ssh danielapi@danielapi.local 'cd ~/services/amelia-porfolio && tar czf - pb_data' \
  > backups/amelia-pb_data-$(date +%F).tar.gz
```

## Updates futuros

```bash
./cmd/build-pi.sh
scp dist/porfolio-amelia danielapi@danielapi.local:~/services/amelia-porfolio/dist/
ssh danielapi@danielapi.local 'cd ~/services/amelia-porfolio && docker compose up -d --build'
```

`pb_data/` no se vuelve a copiar en updates: el volumen lo preserva.
