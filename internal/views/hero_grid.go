package views

import (
	"fmt"
	"html"
	"strings"
)

// renderHeroImagesGrid devuelve HTML con el grid actual del carrusel hero,
// con botones HTMX para reordenar (↑ ↓) y eliminar (×) cada imagen.
// Pensado para usarse desde admin_templ.go y desde el handler HTMX que
// devuelve solo el grid actualizado.
func renderHeroImagesGrid(imgs []AdminImage) string {
	return RenderHeroImagesGrid(imgs)
}

func RenderHeroImagesGrid(imgs []AdminImage) string {
	if len(imgs) == 0 {
		return `<div id="hero-images-grid" class="text-xs text-gray-500 mb-2 p-3 border border-dashed rounded">Sin imágenes en el carrusel.</div>`
	}
	var b strings.Builder
	b.WriteString(`<div id="hero-images-grid" class="grid grid-cols-3 sm:grid-cols-4 gap-2 mb-2">`)
	for i, img := range imgs {
		fn := html.EscapeString(img.Filename)
		url := html.EscapeString(img.URL)
		canUp := i > 0
		canDown := i < len(imgs)-1
		fmt.Fprintf(&b, `
<div class="relative group border border-gray-200 rounded overflow-hidden bg-gray-50">
  <img src="%s" class="w-full h-28 object-cover block" alt="">
  <div class="absolute top-1 left-1 text-[10px] bg-black/60 text-white px-1.5 rounded">%d</div>
  <div class="absolute top-1 right-1 flex gap-1 opacity-0 group-hover:opacity-100 transition">`, url, i+1)
		if canUp {
			fmt.Fprintf(&b, `<button type="button" hx-post="/admin/settings/hero_images/move" hx-vals='{"filename":"%s","dir":"up"}' hx-target="#hero-images-grid" hx-swap="outerHTML" class="bg-white/90 hover:bg-white text-gray-800 w-6 h-6 rounded text-xs leading-none border border-gray-300" title="Subir">↑</button>`, fn)
		}
		if canDown {
			fmt.Fprintf(&b, `<button type="button" hx-post="/admin/settings/hero_images/move" hx-vals='{"filename":"%s","dir":"down"}' hx-target="#hero-images-grid" hx-swap="outerHTML" class="bg-white/90 hover:bg-white text-gray-800 w-6 h-6 rounded text-xs leading-none border border-gray-300" title="Bajar">↓</button>`, fn)
		}
		fmt.Fprintf(&b, `<button type="button" hx-post="/admin/settings/hero_images/delete" hx-vals='{"filename":"%s"}' hx-target="#hero-images-grid" hx-swap="outerHTML" hx-confirm="¿Eliminar esta imagen del carrusel?" class="bg-red-600 hover:bg-red-700 text-white w-6 h-6 rounded text-xs leading-none" title="Eliminar">×</button>`, fn)
		b.WriteString(`</div></div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}
