# Guía del Panel de Administración

Todo lo que se puede editar desde **<https://ameporfolio.danielaregert.com.ar/admin>**.

> Para usuarias y operadores. Si sos agente IA, mirá también `DEPLOY.md`.

---

## Acceder

URL: <https://ameporfolio.danielaregert.com.ar/admin/login>

Credenciales actuales:
- **amelia.repetto@gmail.com** — Amelia
- **admin@ameliarepetto.com** — técnico

Sesión por cookie. Cerrá sesión con **Salir** (arriba a la derecha).

---

## Navegación

Sidebar oscura a la izquierda con cuatro secciones principales + acceso al sitio público.

| Item            | Para qué                                                     |
|-----------------|--------------------------------------------------------------|
| **Secciones**   | Categorías del portfolio (Teatro, Cine, Performance, etc.)   |
| **Obras**       | Cada obra, película, performance, taller, videoclip          |
| **Prensa**      | Notas independientes (no asociadas a una obra)               |
| **Configuración** | Datos personales, bio, contacto, redes, fotos del hero     |
| **Ver sitio**   | Abre el portfolio público en una pestaña nueva               |

En mobile la sidebar se colapsa; tocá ☰ para abrirla.

---

## 1 · Secciones (`/admin`)

Las **categorías principales** del portfolio. Por defecto vienen:
Performance · Teatro · Cine · Videoclips · Talleres · Salud Mental.

### Acciones disponibles

- **Listar todas** con: nombre, slug, orden, cantidad de obras, estado activo/inactivo.
- **Nueva sección** (botón arriba a la derecha):
  - **Nombre** (obligatorio): cómo se ve en el sitio (ej. "Cine").
  - **Slug** (obligatorio): URL amigable (ej. "cine"). Solo letras minúsculas, números y guiones.
  - **Descripción**: texto corto opcional.
  - **Orden**: número que define en qué posición sale (menor = primero).
  - **Activa**: tildá para que aparezca en el sitio.
  - **Imagen de portada**: opcional, aparece si el frontend la usa.
- **Editar** una existente (clic en su nombre).
- **Activar / desactivar**: con el toggle del listado, no se borra; solo se oculta.
- **Eliminar**: borra la sección **y todas las obras adentro**. Pide confirmación.

> Tip: cambiar el **slug** de una sección rompe links externos a esa categoría.
> Si publicaste el link en redes, no lo cambies sin razón.

---

## 2 · Obras (`/admin/works`)

El corazón del portfolio: una entrada por cada obra de teatro, performance,
película, taller, videoclip, etc.

### Listado

Vista con: thumbnail, título, año, sección, rol, estado, cantidad de links,
botón ⭐ destacada.

Filtros: por sección (selector arriba).

### Crear / Editar una obra

Campos del formulario:

| Campo           | Notas                                                         |
|-----------------|---------------------------------------------------------------|
| **Sección**     | A qué categoría pertenece (relación obligatoria).             |
| **Título**      | Obligatorio.                                                  |
| **Slug**        | URL amigable, único. Ej: `bzd-el-entierro`.                   |
| **Año**         | Texto libre. Acepta "2024", "2019-2021", "2024 (en gira)".    |
| **Rol**         | Tu rol en la obra: "Actriz", "Dramaturgia y dirección", etc.  |
| **Descripción** | Texto largo. Acepta HTML básico (`<p>`, `<em>`, `<a>`, etc.). |
| **Créditos**    | Lista de equipo: "Dirección: X · Dramaturgia: Y · ..."        |
| **Orden**       | Número, menor = más arriba en la sección.                     |
| **Activa**      | Tildá para mostrarla en el sitio.                             |
| **Destacada**   | Tildá para resaltarla (depende del frontend).                 |

### Multimedia de la obra

Tres tipos de archivos por obra:

#### Imágenes (hasta 12)

- Input de subida múltiple: ⌘-click / ctrl-click para elegir varias.
- Las imágenes nuevas se **agregan** al final.
- Cada thumbnail tiene una **× roja** al hover para eliminarla individualmente.
- Sin reorder por ahora desde la UI (roadmap).
- Formatos aceptados: JPG, PNG, WebP, AVIF. Máx. 10MB por imagen.

#### Video (opcional, 1 archivo)

- Input "Video (mp4 / mov, opcional)".
- Si subís uno nuevo, reemplaza el anterior.
- Para borrar el actual sin reemplazar: tildá **"borrar video actual"** y guardá.
- Formatos: mp4, webm, mov (quicktime). Máx. 200MB.
- Aparece embebido en el sitio público con controles HTML5.

#### Dossier (PDF, 1 archivo)

- Input "Dossier (PDF, opcional)".
- Reemplazo o borrado igual que el video.
- Máx. 50MB.
- Aparece como botón "Descargar dossier (PDF)" junto a los links.

### Links de la obra

Sección debajo del form principal (solo visible en obras ya creadas).

- **+ Agregar link** abre un mini-form: Label, URL, Tipo.
- Tipos disponibles: `trailer` · `full_work` (obra completa) · `dossier` · `press` · `external` · `instagram`.
- **Eliminar link**: botón rojo a la derecha de cada uno.
- Sin reorder por ahora desde la UI.

### Eliminar una obra

Botón "Eliminar" abajo del form, con confirmación. **Borra también todos sus
links e imágenes asociadas.**

---

## 3 · Prensa (`/admin/press`)

Notas que **no** están vinculadas a una obra específica (cobertura general,
entrevistas, perfiles).

### Acciones

- Listar / crear / editar / borrar / activar-desactivar.
- Campos:
  - **Título** (obligatorio)
  - **Publicación**: medio (ej. "El Mundo")
  - **URL**: link al artículo
  - **Fecha**: texto libre, ej. "2023-09-04"
  - **Extracto**: párrafo o cita corta. Acepta HTML básico.
  - **Orden** y **Activa**.

> Diferencia con los **links de prensa** dentro de una obra: si la nota habla
> específicamente de una obra, agregalo como link de esa obra (kind: `press`).
> Si es una nota general sobre Amelia, va acá.

---

## 4 · Configuración (`/admin/settings`)

Página única con todo lo del sitio. Cambios afectan al instante el frontend.

### Bloque "General"

- **Site name**: nombre que aparece en el `<title>` del navegador y en el nav.
- **Tagline**: bajada (ej. "porfolio artístico").

### Bloque "Biografía"

- **Bio (Español)**: bio completa para la home pública. Acepta HTML
  (`<p>`, `<em>`, `<a href="...">`, `<strong>`, etc.).
- **Bio (English)**: versión en inglés (no se muestra todavía en el sitio,
  reservado para una futura versión bilingüe).

### Bloque "Contacto"

- **Email**: aparece como link `mailto:` en el sitio.
- **Móvil**.
- **Dirección**.
- **☑ Mostrar la dirección en la web pública**: si lo destildás, la dirección
  desaparece del sitio (queda guardada igual). Útil por privacidad.

### Bloque "Redes sociales"

URLs completas (con `https://`). Si dejás vacío, no se muestra el link.
- Instagram · Facebook · YouTube · Vimeo

### Bloque "Imagen / Reel"

- **Profile image**: foto de perfil (1 archivo).
- **Hero image**: imagen principal del hero si **no** se usa el carrusel
  (fallback).
- **Reel URL** (texto): URL del reel si querés mostrarlo.

### Bloque "Carrusel hero" — el más interactivo

Hasta **8 imágenes** que rotan automáticamente con el efecto glitch en el hero.

#### Grid de imágenes actuales

Aparece arriba del input de subida. Cada thumbnail muestra su **número de
orden** (1, 2, 3...) y al hacer hover aparecen 3 botones:

| Botón | Acción                                                  |
|-------|---------------------------------------------------------|
| **↑** | Sube esa imagen un puesto en el orden                   |
| **↓** | La baja un puesto                                       |
| **×** | La elimina (pide confirmación)                          |

Los cambios se aplican al instante (HTMX) sin recargar la página y se
reflejan en el sitio público inmediatamente.

#### Subir más imágenes

Input "Carrusel hero (hasta 8 imágenes)" → seleccioná varias con
⌘-click / ctrl-click.

**Importante:** subir nuevas imágenes las **agrega al final** del carrusel
(no reemplaza). Si querés vaciar el carrusel, eliminalas una por una con la ×.

> Tip de diseño: las imágenes con proporción **vertical (4:5)** funcionan
> mejor porque el frontend usa `aspect-ratio: 4/5` con `object-fit: cover`.
> Imágenes muy panorámicas se recortan en alto.

### Botón **Guardar configuración**

Aplica los cambios del formulario. El grid del carrusel ya guarda solo
(porque usa HTMX), pero el resto requiere apretar este botón.

---

## Cosas que NO se editan desde el admin

| Para cambiar...                      | Necesitás...                              |
|--------------------------------------|-------------------------------------------|
| El layout del sitio (colores, fonts) | Editar `internal/handlers/public.go` + redeploy |
| Nuevos tipos de campo en obras       | Crear migración (`migrations/0NN_*.go`) + redeploy |
| Hover colors del navbar              | Editar CSS en `public.go`                 |
| El efecto glitch                     | Editar CSS en `public.go`                 |
| Crear un nuevo superuser             | Comando en el Pi (ver `DEPLOY.md`)        |

---

## FAQ rápido

**¿Cómo agrego una obra nueva?**
1. Andá a Obras → "Nueva obra".
2. Elegí la sección, completá título, slug, año, rol.
3. Guardá.
4. En la pantalla de edición, subí imágenes / video / dossier y agregá links.

**¿Por qué mi obra no aparece en el sitio?**
- ¿La marcaste como **Activa**?
- ¿La sección a la que pertenece está activa?
- Refresca con `Ctrl+Shift+R` (Cloudflare puede cachear).

**Borré algo por error.**
- Si fue hace poco y tenemos backup automático, se restaura desde el último
  snapshot (ver `DEPLOY.md` → Backups → Restaurar).
- Borrar es destructivo. **Desactivar** (toggle) es reversible y suele ser
  la mejor opción si dudás.

**¿Cuánto tardan los cambios en verse en el sitio?**
- Inmediato. Si usás Cloudflare con cache agresivo, podés tener que esperar
  o purgar.

**¿Puedo subir un MOV grande?**
- Hasta 200MB. Si pesa más, conviene comprimirlo (Handbrake, `ffmpeg`) a
  H.264 mp4 antes de subir.
