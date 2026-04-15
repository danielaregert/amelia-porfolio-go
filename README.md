# Porfolio Amelia Repetto

Portfolio artístico bilingüe (ES/EN) con CMS propio. Stack minimalista en Go
+ PocketBase + Tailwind, deployado en Raspberry Pi 5 atrás de Cloudflare Tunnel.

**En producción:** <https://ameporfolio.danielaregert.com.ar>

---

## Para qué sirve cada cosa

| Doc                          | Para qué                                                 |
|------------------------------|----------------------------------------------------------|
| **README.md** (este archivo) | Inicio rápido, stack, desarrollo local, comandos claves. |
| [`ADMIN.md`](./ADMIN.md)     | Guía del panel de administración (para Amelia + operadores). |
| [`DEPLOY.md`](./DEPLOY.md)   | Operaciones: deploy, backups, troubleshooting, updates.  |

---

## Stack

- **Backend:** Go 1.26 + [PocketBase](https://pocketbase.io) v0.36 (embebido como librería).
- **DB:** SQLite (en `pb_data/data.db`), con WAL.
- **Frontend público:** HTML servido por Go (`html/template`) con Tailwind CDN,
  Roboto Slab (variable font), efectos CSS glitch y lightbox HTML5 vanilla.
- **Admin:** [templ](https://templ.guide) + Tailwind + HTMX + Alpine.js, estilo
  TailAdmin (sidebar oscura).
- **i18n:** ES default, EN en `/en/`. Toggle en navbar, `hreflang`, og:locale.
- **Infra prod:** Docker Compose en Raspberry Pi 5, Caddy como reverse proxy
  (red `web`), Cloudflare Tunnel para TLS + DNS público.

---

## Estructura del repo

```
.
├── main.go                      entrypoint PocketBase + registro de routes
├── cmd/build-pi.sh              cross-compile linux/arm64 para el Pi
├── deploy.sh                    script orquestador del deploy (ver DEPLOY.md)
├── Dockerfile                   runtime slim con el binario
├── docker-compose.yml           service amelia-porfolio en red Docker "web"
│
├── internal/
│   ├── handlers/
│   │   ├── routes.go            registro de rutas (públicas + admin)
│   │   ├── public.go            frontend (HTML inline, i18n, hero glitch,
│   │   │                         galería lightbox, meta/OG tags)
│   │   └── admin.go             CRUD, sesiones cookie, rate-limit login
│   ├── views/
│   │   ├── admin.templ          SOURCE del admin UI (editar este, no el .go)
│   │   ├── admin_templ.go       GENERADO por `templ generate`
│   │   ├── types.go             structs de formularios del admin
│   │   └── hero_grid.go         fragment HTMX para grid del carrusel hero
│   ├── adminsession/            auth por cookie (HttpOnly + SameSite)
│   ├── ratelimit/               rate-limit del login
│   └── sanitize/                bluemonday wrapper para HTML de bios/descripciones
│
├── migrations/                  schema + seed + cargas de imágenes
│   ├── 001_portfolio_schema.go
│   ├── 002_seed_amelia.go
│   ├── 003-014_*.go             cargas, fixes, carrusel, i18n
│   └── 015-019_*.go             schema i18n, traducciones, meta OG
│
├── static/                      archivos estáticos embebidos (favicon, etc)
│
├── pb_data/                     ⚠️ NO COMMITEADO (gitignored)
│   ├── data.db                  SQLite (tablas + config + migraciones aplicadas)
│   ├── auxiliary.db             logs y sesiones
│   └── storage/                 uploads del admin (imágenes, videos, PDFs)
│
├── ADMIN.md                     guía del panel de admin
├── DEPLOY.md                    guía de deploy y operaciones
└── README.md                    este archivo
```

---

## Pre-requisitos

### Desarrollo local (Mac / Linux)
- **Go 1.26+** (`brew install go`)
- **templ** CLI (`go install github.com/a-h/templ/cmd/templ@latest` → queda en `~/go/bin/templ`)
- **sqlite3** CLI (viene con macOS)

### Deploy
- SSH con clave autorizada al Pi (`ssh danielaregert@danielapi.local` sin password)
- Caddy en Docker ya corriendo en el Pi, conectado a la red `web`
- Cloudflare Tunnel ya configurado apuntando al dominio

---

## Inicio rápido

### 1. Clonar y traer datos

```bash
git clone https://github.com/danielaregert/amelia-porfolio-go.git
cd amelia-porfolio-go

# Traer el estado real de producción (DB + storage)
rsync -azh danielaregert@danielapi.local:~/services/amelia-porfolio/pb_data/ pb_data/
```

### 2. Correr local

```bash
go run . serve --http=127.0.0.1:8099
```

Abrir:
- <http://127.0.0.1:8099/> — portfolio público
- <http://127.0.0.1:8099/en/> — versión inglés
- <http://127.0.0.1:8099/admin/login> — panel de admin
  - Credenciales ver Bitwarden o `DEPLOY.md`

### 3. Primer superusuario (si la DB está vacía)

```bash
go run . superuser create tu@email.com "PASSWORD"
```

---

## Workflow de cambios

### Cambio de código (frontend, handler, CSS, JS)

```bash
# 1. Editar archivos
# 2. Probar local
go run . serve --http=127.0.0.1:8099
# Revisar en el browser.

# 3. Commit
git add -A
git commit -m "describe cambio"
git push

# 4. Deploy
./deploy.sh
```

### Cambio de admin UI (archivos `.templ`)

Los archivos `internal/views/admin.templ` son el **source**. El `_templ.go` se
genera automáticamente. Después de editar:

```bash
~/go/bin/templ generate   # regenera admin_templ.go
go run . serve --http=127.0.0.1:8099   # probar
./deploy.sh               # el build-pi.sh también corre templ generate
```

### Cambio de schema (nueva migración)

```bash
# Crear migrations/0NN_descripcion.go
# Probar local: automigrate corre al arrancar
go run . serve --http=127.0.0.1:8099
sqlite3 pb_data/data.db "..."   # verificar

# Deploy con backup (por si algo sale mal)
./deploy.sh --backup --logs
```

**Convenciones de migraciones:**
- Numeradas `0NN_descripcion.go`, arranca en 001.
- Idempotentes cuando sea posible: `if col.Fields.GetByName("X") == nil { add }`.
- La función `down` puede ser `func(app core.App) error { return nil }` si no
  vale la pena revertir.
- `Automigrate: true` en `main.go` las aplica al arrancar.

### Cambio de contenido (textos, imágenes, videos, links)

**No requiere deploy.** Amelia u operador entra al panel:

<https://ameporfolio.danielaregert.com.ar/admin/login>

Todo persiste en `pb_data/` del Pi y sobrevive a redeploys.

---

## Comandos útiles

```bash
# Correr local
go run . serve --http=127.0.0.1:8099

# Build para el Pi (manual)
./cmd/build-pi.sh
# → dist/porfolio-amelia (arm64, ~33MB)

# Deploy completo (build + scp + restart + smoke test)
./deploy.sh

# Deploy con backup automático + seguir logs
./deploy.sh --backup --logs

# Re-skip build (cuando solo cambió Dockerfile)
./deploy.sh --no-build

# Ver logs en vivo del container del Pi
ssh danielaregert@danielapi.local 'docker logs -f amelia-porfolio'

# Consultar DB de producción
ssh danielaregert@danielapi.local \
  'docker exec amelia-porfolio sqlite3 /app/pb_data/data.db ".tables"'

# API REST pública (pocketbase)
curl https://ameporfolio.danielaregert.com.ar/api/collections/works/records
```

---

## Dominios y accesos

| URL                                            | Para qué                       |
|------------------------------------------------|--------------------------------|
| <https://ameporfolio.danielaregert.com.ar/>    | Sitio público                  |
| <https://ameporfolio.danielaregert.com.ar/en/> | Sitio público en inglés        |
| <https://ameporfolio.danielaregert.com.ar/admin> | Admin custom                 |
| <https://ameporfolio.danielaregert.com.ar/_/>  | Admin nativo de PocketBase     |
| <https://ameporfolio.danielaregert.com.ar/api/> | REST API pública             |

---

## Checklist antes de shippear cambios importantes

- [ ] `go build ./...` sin errores.
- [ ] `go run . serve --http=127.0.0.1:8099` arranca limpio.
- [ ] Probé el cambio en `/` y `/en/` si afecta el frontend.
- [ ] Probé el cambio en `/admin` si afecta el panel.
- [ ] `git status` no deja archivos basura.
- [ ] `./deploy.sh --backup` si hay migración nueva.
- [ ] Visité el sitio en producción tras el deploy.

---

## Licencia y créditos

Código propietario. Contenido (textos, fotos, videos) propiedad de Amelia Repetto.

Librerías third-party bajo sus licencias respectivas (MIT/Apache-2.0 en su mayoría):
- [PocketBase](https://pocketbase.io) · MIT
- [templ](https://templ.guide) · MIT
- [bluemonday](https://github.com/microcosm-cc/bluemonday) · BSD-3
- [Tailwind CSS](https://tailwindcss.com) · MIT
- [Alpine.js](https://alpinejs.dev) · MIT
- [HTMX](https://htmx.org) · BSD-2
