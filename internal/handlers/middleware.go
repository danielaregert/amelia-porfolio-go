package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

// securityHeaders aplica cabeceras HTTP de seguridad comunes antes de delegar
// en el handler. Defensa en profundidad por si Caddy no las agrega.
func securityHeaders(h func(e *core.RequestEvent) error) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		hdr := e.Response.Header()
		hdr.Set("X-Frame-Options", "DENY")
		hdr.Set("X-Content-Type-Options", "nosniff")
		hdr.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		hdr.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		return h(e)
	}
}

// sameOrigin verifica que la petición POST venga del mismo host. Combinado con
// SameSite=Lax en la cookie de sesión, da una defensa sólida contra CSRF sin
// necesitar tokens por formulario.
func sameOrigin(r *http.Request) bool {
	if r.Method != http.MethodPost {
		return true
	}
	expected := r.Host
	if fwd := r.Header.Get("X-Forwarded-Host"); fwd != "" {
		expected = fwd
	}
	check := func(raw string) bool {
		if raw == "" {
			return false
		}
		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			return false
		}
		return strings.EqualFold(u.Host, expected)
	}
	if o := r.Header.Get("Origin"); o != "" {
		return check(o)
	}
	if ref := r.Header.Get("Referer"); ref != "" {
		return check(ref)
	}
	// Sin Origin ni Referer en un POST: bloquear.
	return false
}

// csrfGuard wraps POST handlers validando el origen. Los handlers GET quedan
// sin alterar.
func csrfGuard(h func(e *core.RequestEvent) error) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if !sameOrigin(e.Request) {
			return e.ForbiddenError("origen inválido", nil)
		}
		return h(e)
	}
}
