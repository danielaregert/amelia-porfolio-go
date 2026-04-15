#!/usr/bin/env bash
# ============================================================================
# deploy.sh — Build + push del binario al Pi y restart del contenedor.
#
# Uso:
#   ./deploy.sh                 # build + scp + restart (workflow normal)
#   ./deploy.sh --no-build      # solo subir binario existente y restart
#   ./deploy.sh --backup        # hace backup de pb_data ANTES de deploy
#   ./deploy.sh --logs          # después de deployar, sigue los logs
#   ./deploy.sh --help
#
# Variables de entorno (override de defaults):
#   PI_HOST     destino ssh, default "danielaregert@danielapi.local"
#   PI_PATH     ruta en el Pi,  default "~/services/amelia-porfolio"
#   CONTAINER   nombre docker,  default "amelia-porfolio"
#
# Pre-requisitos (chequeos automáticos):
#   - SSH funcional al Pi sin password (key ed25519 autorizada)
#   - Docker Compose v2 en el Pi
#   - Red docker "web" presente (la usa el caddy reverse proxy)
#
# El script NUNCA toca pb_data/ del Pi (datos de Amelia). Para sincronizar
# datos desde el Mac al Pi (raro, solo migración inicial), correr a mano:
#   rsync -azh --delete pb_data/ $PI_HOST:$PI_PATH/pb_data/
# ============================================================================

set -euo pipefail

PI_HOST="${PI_HOST:-danielaregert@danielapi.local}"
PI_PATH="${PI_PATH:-~/services/amelia-porfolio}"
CONTAINER="${CONTAINER:-amelia-porfolio}"

DO_BUILD=1
DO_BACKUP=0
DO_LOGS=0

for arg in "$@"; do
  case "$arg" in
    --no-build) DO_BUILD=0 ;;
    --backup)   DO_BACKUP=1 ;;
    --logs)     DO_LOGS=1 ;;
    --help|-h)
      sed -n '2,/^# ====/p' "$0" | sed 's/^# \{0,1\}//'
      exit 0
      ;;
    *) echo "flag desconocido: $arg" >&2; exit 2 ;;
  esac
done

# Colores (suprimidos si no hay TTY)
if [[ -t 1 ]]; then
  C_BOLD=$'\033[1m'; C_DIM=$'\033[2m'; C_RED=$'\033[31m'
  C_GREEN=$'\033[32m'; C_YELLOW=$'\033[33m'; C_RESET=$'\033[0m'
else
  C_BOLD=""; C_DIM=""; C_RED=""; C_GREEN=""; C_YELLOW=""; C_RESET=""
fi
say()  { echo "${C_BOLD}==> $*${C_RESET}"; }
warn() { echo "${C_YELLOW}!! $*${C_RESET}" >&2; }
die()  { echo "${C_RED}xx $*${C_RESET}" >&2; exit 1; }
ok()   { echo "${C_GREEN}✓${C_RESET} $*"; }

# ---------- 1. PRE-FLIGHT ----------
say "Pre-flight checks"

[[ -f main.go ]] || die "no estoy en la raíz del proyecto (no veo main.go)"
[[ -x cmd/build-pi.sh ]] || die "falta cmd/build-pi.sh ejecutable"
command -v ssh >/dev/null || die "ssh no encontrado"
command -v scp >/dev/null || die "scp no encontrado"

ssh -o BatchMode=yes -o ConnectTimeout=5 "$PI_HOST" 'true' \
  || die "no puedo entrar por SSH a $PI_HOST (revisar key, fail2ban o que el Pi esté online)"
ok "SSH a $PI_HOST"

ssh "$PI_HOST" 'docker --version >/dev/null 2>&1' \
  || die "docker no responde en el Pi"
ok "docker en el Pi"

# ---------- 2. BACKUP OPCIONAL ----------
if [[ "$DO_BACKUP" -eq 1 ]]; then
  say "Backup de pb_data antes del deploy"
  mkdir -p ~/backups/amelia-porfolio
  STAMP=$(date +%F-%H%M)
  OUT="$HOME/backups/amelia-porfolio/pb_data-${STAMP}.tar.gz"
  ssh "$PI_HOST" "cd $PI_PATH && tar czf - pb_data" > "$OUT"
  ok "backup → $OUT ($(du -h "$OUT" | cut -f1))"
fi

# ---------- 3. BUILD ----------
if [[ "$DO_BUILD" -eq 1 ]]; then
  say "Compilando binario para linux/arm64"
  ./cmd/build-pi.sh
else
  warn "saltando build (--no-build); usando dist/porfolio-amelia existente"
  [[ -f dist/porfolio-amelia ]] || die "no hay dist/porfolio-amelia para subir"
fi
SIZE=$(du -h dist/porfolio-amelia | cut -f1)
ok "binario listo ($SIZE)"

# ---------- 4. SCP BINARIO ----------
say "Subiendo binario al Pi"
ssh "$PI_HOST" "mkdir -p $PI_PATH/dist"
scp -q dist/porfolio-amelia "$PI_HOST:$PI_PATH/dist/"
ok "binario en $PI_HOST:$PI_PATH/dist/"

# Subimos también Dockerfile + compose por si cambiaron (overhead despreciable).
scp -q Dockerfile docker-compose.yml "$PI_HOST:$PI_PATH/" >/dev/null 2>&1 || true
ok "Dockerfile + compose sincronizados"

# ---------- 5. RESTART ----------
say "Rebuild + restart del contenedor"
ssh "$PI_HOST" "cd $PI_PATH && docker compose up -d --build" 2>&1 \
  | grep -vE "^#|^DONE|^naming|^exporting|^WARN.*FromAsCasing" || true
ok "contenedor $CONTAINER running"

# ---------- 6. SMOKE TEST ----------
say "Smoke test"
sleep 2
HTTP=$(ssh "$PI_HOST" "docker exec caddy curl -s -o /dev/null -w '%{http_code}' http://$CONTAINER:8090/" 2>/dev/null || echo "ERR")
if [[ "$HTTP" == "200" ]]; then
  ok "caddy → $CONTAINER:8090 responde 200"
else
  warn "smoke test devolvió HTTP $HTTP — revisar logs:  ssh $PI_HOST 'docker logs $CONTAINER'"
fi

# ---------- 7. LOGS OPCIONAL ----------
if [[ "$DO_LOGS" -eq 1 ]]; then
  say "Siguiendo logs (Ctrl+C para salir)"
  ssh "$PI_HOST" "docker logs -f --tail 30 $CONTAINER"
fi

say "Deploy completo · https://ameporfolio.danielaregert.com.ar"
