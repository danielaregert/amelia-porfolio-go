# Deploy & Operations Guide

> Guía rápida para humanos y agentes (Claude Code u otros) que necesiten
> deployar, restaurar, o modificar el portfolio de Amelia Repetto.
>
> **Pila:** Go + PocketBase, Docker en Raspberry Pi 5, Caddy reverse proxy,
> Cloudflare Tunnel para TLS, dominio `ameporfolio.danielaregert.com.ar`.

---

## TL;DR

```bash
# Edito código → pruebo local → deploy
go run . serve --http=127.0.0.1:8099   # local
./deploy.sh                            # producción
```

Eso es todo el ciclo. El script chequea todo y te avisa si algo falla.

---

## Mapa mental de la app

| Pieza                 | Dónde vive                                              | Cuándo cambia |
|-----------------------|---------------------------------------------------------|---------------|
| Código fuente         | GitHub `danielaregert/amelia-porfolio-go` + Mac local   | cuando vos editás |
| Binario `arm64`       | `dist/porfolio-amelia` (Mac) y `~/services/amelia-porfolio/dist/` (Pi) | cada deploy |
| Schema DB + relaciones| Migraciones en `migrations/*.go`                        | cuando agregás migración |
| Datos (textos, links) | `pb_data/data.db` (SQLite) **en el Pi**                 | cuando Amelia edita en admin |
| Imágenes / video / PDF| `pb_data/storage/...` **en el Pi**                      | cuando Amelia sube por admin |
| TLS / DNS             | Cloudflare (Tunnel + cert)                              | rara vez |

**Regla de oro:** `pb_data/` del Pi es **sagrado**. Nunca lo sobreescribas en
deploys normales. Solo se toca para:
- Backups (lectura)
- Restore de un backup (escritura controlada)
- Bootstrap inicial (sync una vez del Mac)

---

## El script `deploy.sh`

### Qué hace, en orden

1. **Pre-flight**: verifica que estás en la raíz del repo, que tenés SSH al Pi
   sin password, y que Docker responde en el Pi.
2. **Backup opcional** (`--backup`): tar de `pb_data/` del Pi → `~/backups/amelia-porfolio/pb_data-FECHA.tar.gz` en tu Mac.
3. **Build**: corre `cmd/build-pi.sh` que cross-compila `linux/arm64` →
   `dist/porfolio-amelia`. Saltable con `--no-build`.
4. **SCP**: sube el binario + Dockerfile + docker-compose.yml al Pi.
5. **Restart**: en el Pi corre `docker compose up -d --build`. Reusa la
   imagen base si no hay cambios en `Dockerfile` (rebuild rapidísimo, solo
   reemplaza la capa del binario).
6. **Smoke test**: pide a Caddy (que está en la misma red Docker `web`) que
   haga `curl http://amelia-porfolio:8090/` y espera `HTTP 200`.
7. **Logs opcionales** (`--logs`): tail de los logs del contenedor.

### Flags

| Flag         | Efecto                                                |
|--------------|-------------------------------------------------------|
| `--no-build` | Saltea la compilación, usa el binario actual          |
| `--backup`   | Backup de `pb_data/` ANTES del deploy                 |
| `--logs`     | Sigue logs después del deploy (Ctrl+C para salir)     |
| `--help`     | Muestra ayuda                                         |

### Variables de entorno (overrides)

```bash
PI_HOST=otrouser@otropi.local ./deploy.sh    # otro destino SSH
PI_PATH=/opt/amelia ./deploy.sh              # otro path en el Pi
CONTAINER=amelia-staging ./deploy.sh         # otro nombre de contenedor
```

### Ejemplos típicos

```bash
./deploy.sh                       # workflow normal: edito, deploy
./deploy.sh --backup              # backup primero (hacelo antes de migraciones)
./deploy.sh --no-build --logs     # cambié solo Dockerfile o quiero ver logs
./deploy.sh --backup --logs       # paranoia + observación
```

---

## Workflows

### Cambio de código (handler, frontend, CSS)

```bash
# Edito archivos en el Mac, pruebo:
go run . serve --http=127.0.0.1:8099
# Curl, navegador, lo que sea.

# Cuando está listo:
git add -A && git commit -m "describe el cambio"
git push
./deploy.sh
```

### Cambio de schema o seed (nueva migración)

```bash
# Crear archivo migrations/0NN_descripcion.go
# Probar local primero — automigrate corre al arrancar `go run`
go run . serve --http=127.0.0.1:8099
# Verificar con: sqlite3 pb_data/data.db "..."

# Backup obligatorio antes de tocar producción:
git add -A && git commit -m "migration NN: ..."
git push
./deploy.sh --backup --logs
# Mirá los logs para confirmar que la migración aplicó OK.
```

### Cambio de contenido (textos, fotos, links)

**No requiere deploy.** Amelia edita en admin:

- URL: <https://ameporfolio.danielaregert.com.ar/admin/login>
- Cambios persisten en `pb_data/` del Pi y sobreviven a redeploys.

---

## Backups

### Manual desde Mac

```bash
mkdir -p ~/backups/amelia-porfolio
ssh danielaregert@danielapi.local 'cd ~/services/amelia-porfolio && tar czf - pb_data' \
  > ~/backups/amelia-porfolio/pb_data-$(date +%F-%H%M).tar.gz
```

### Automático en el Pi (cron semanal)

```bash
ssh danielaregert@danielapi.local
mkdir -p ~/backups/amelia-porfolio
crontab -e
# Pegar:
0 3 * * 0 cd ~/services/amelia-porfolio && tar czf ~/backups/amelia-porfolio/pb_data-$(date +\%F).tar.gz pb_data && ls -t ~/backups/amelia-porfolio/*.tar.gz | tail -n +9 | xargs -r rm
```

Domingos 3am, mantiene los últimos 8 snapshots.

### Restaurar

```bash
ssh danielaregert@danielapi.local
cd ~/services/amelia-porfolio
docker compose down
mv pb_data pb_data.bad
tar xzf ~/backups/amelia-porfolio/pb_data-2026-04-20.tar.gz   # ajustá fecha
docker compose up -d
docker logs -f amelia-porfolio
```

---

## Crear / resetear superusers

```bash
# Crear nuevo superuser (admin de PocketBase)
ssh danielaregert@danielapi.local \
  'docker exec amelia-porfolio /app/porfolio-amelia superuser create EMAIL "PASSWORD"'

# Resetear password de uno existente
ssh danielaregert@danielapi.local \
  'docker exec amelia-porfolio /app/porfolio-amelia superuser update EMAIL "NUEVA"'
```

Superusers actuales:
- `admin@ameliarepetto.com` (técnico)
- `amelia.repetto@gmail.com` (de Amelia)

---

## Infra: dónde está cada cosa

```
Internet ─► Cloudflare (TLS + Tunnel)
            │
            ▼
       Pi en LAN ─► Caddy (Docker, red "web", :80)
                    │ reverse_proxy http://amelia-porfolio:8090
                    ▼
              amelia-porfolio (Docker, red "web")
                    │ binario Go + PocketBase
                    │ volumen → ./pb_data
                    ▼
              ~/services/amelia-porfolio/pb_data/
                    ├── data.db          (SQLite)
                    ├── auxiliary.db     (logs/sesiones)
                    └── storage/         (uploads)
```

### Archivos clave del repo

| Archivo                      | Para qué                                          |
|------------------------------|---------------------------------------------------|
| `main.go`                    | entrypoint PocketBase + registro de routes        |
| `internal/handlers/public.go`| frontend público (HTML inline + Roboto Slab + glitch) |
| `internal/handlers/admin.go` | admin con sesiones cookie + rate-limit            |
| `internal/views/admin_templ.go` | layout admin (TailAdmin-style, sidebar oscura)  |
| `migrations/0NN_*.go`        | schema + seed (corren al arrancar)                |
| `Dockerfile`                 | runtime slim (debian-bookworm + binario)          |
| `docker-compose.yml`         | service `amelia-porfolio` en red `web`            |
| `cmd/build-pi.sh`            | cross-compile linux/arm64                         |
| `deploy.sh`                  | el orquestador del deploy                         |
| `deploy/Caddyfile.snippet`   | bloque para tu Caddyfile global del Pi            |

---

## Troubleshooting

### `ssh: Connection refused`
- Pi viva: `ping danielapi.local`
- Si responde pero :22 rechaza → fail2ban te baneó. Entrá por Cockpit
  (https://192.168.0.29:9090) → Terminal:
  ```
  sudo fail2ban-client unban --all
  ```

### `docker compose up` falla con "network web not found"
La red Docker que usa Caddy debe existir antes:
```bash
ssh danielaregert@danielapi.local 'docker network ls | grep web'
# si no está: docker network create web
# y reconectar el caddy a ella
```

### El sitio devuelve 502 / página no carga
```bash
ssh danielaregert@danielapi.local 'docker logs --tail 50 amelia-porfolio'
ssh danielaregert@danielapi.local 'docker logs --tail 50 caddy'
```

### Una migración rompió la DB
1. `docker compose down` en el Pi
2. Restaurar último backup (ver sección Backups → Restaurar)
3. Investigar la migración en local antes de re-deployar

### Después de deploy el browser sigue mostrando lo viejo
Cloudflare cachea. Forzá:
```bash
curl -X POST -H 'CF-API-Token: ...' \
  'https://api.cloudflare.com/client/v4/zones/ZONE_ID/purge_cache' \
  -d '{"purge_everything":true}'
```
O en el dashboard CF → Caching → Purge Everything. Para evitar caché en
contenido dinámico, agregar Cache Rules en CF.

---

## Para agentes (Claude Code, etc.)

Si te toca operar este proyecto:

1. **Antes de deployar cambios de schema/seed**: corré `./deploy.sh --backup --logs`.
2. **Nunca hagas** `rsync ... pb_data/` en sentido Mac→Pi salvo bootstrap inicial.
3. **Para introspección de datos en producción**:
   ```bash
   ssh danielaregert@danielapi.local \
     'docker exec amelia-porfolio sqlite3 /app/pb_data/data.db "SELECT ..."'
   ```
4. **Para correr una migración manualmente** (raro, automigrate ya lo hace):
   ```bash
   ssh danielaregert@danielapi.local \
     'docker exec amelia-porfolio /app/porfolio-amelia migrate up'
   ```
5. **Logs en vivo**: `ssh ... 'docker logs -f amelia-porfolio'`.
6. **Convención de migraciones**: numeradas `0NN_descripcion.go`, idempotentes
   cuando sea posible (usar `if col.Fields.GetByName("X") == nil { ... }`).
   La función `down` puede ser `func(app core.App) error { return nil }` si el
   cambio no se quiere revertir automáticamente.
