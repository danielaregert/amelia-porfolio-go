package handlers

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func RegisterRoutes(r *router.Router[*core.RequestEvent], app *pocketbase.PocketBase) {
	// Helpers:
	// - pub: GET público con cabeceras de seguridad.
	// - adm: GET admin (headers + requireAdmin dentro del handler).
	// - post: POST con headers + validación de origen (CSRF guard).
	pub := func(h func(e *core.RequestEvent) error) func(e *core.RequestEvent) error {
		return securityHeaders(h)
	}
	post := func(h func(e *core.RequestEvent) error) func(e *core.RequestEvent) error {
		return securityHeaders(csrfGuard(h))
	}

	// --- Public (Phase 2) ---
	r.GET("/", pub(func(e *core.RequestEvent) error { return publicHome(e, app) }))
	r.GET("/en", pub(func(e *core.RequestEvent) error { return publicHome(e, app) }))
	r.GET("/en/", pub(func(e *core.RequestEvent) error { return publicHome(e, app) }))
	r.GET("/robots.txt", pub(func(e *core.RequestEvent) error { return publicRobots(e) }))

	// --- Admin: auth ---
	r.GET("/admin/login", pub(func(e *core.RequestEvent) error { return adminLoginView(e, app) }))
	r.POST("/admin/login", post(func(e *core.RequestEvent) error { return adminLoginSubmit(e, app) }))
	r.POST("/admin/logout", post(func(e *core.RequestEvent) error { return adminLogout(e, app) }))

	// --- Admin: sections ---
	r.GET("/admin", pub(func(e *core.RequestEvent) error { return adminSectionsList(e, app) }))
	r.GET("/admin/sections/new", pub(func(e *core.RequestEvent) error { return adminSectionNewView(e, app) }))
	r.POST("/admin/sections/new", post(func(e *core.RequestEvent) error { return adminSectionNewSubmit(e, app) }))
	r.GET("/admin/sections/{id}", pub(func(e *core.RequestEvent) error { return adminSectionEditView(e, app) }))
	r.POST("/admin/sections/{id}", post(func(e *core.RequestEvent) error { return adminSectionEditSubmit(e, app) }))
	r.POST("/admin/sections/{id}/toggle", post(func(e *core.RequestEvent) error { return adminSectionToggle(e, app) }))
	r.POST("/admin/sections/{id}/delete", post(func(e *core.RequestEvent) error { return adminSectionDelete(e, app) }))

	// --- Admin: works ---
	r.GET("/admin/works", pub(func(e *core.RequestEvent) error { return adminWorksList(e, app) }))
	r.GET("/admin/works/new", pub(func(e *core.RequestEvent) error { return adminWorkNewView(e, app) }))
	r.POST("/admin/works/new", post(func(e *core.RequestEvent) error { return adminWorkNewSubmit(e, app) }))
	r.GET("/admin/works/{id}", pub(func(e *core.RequestEvent) error { return adminWorkEditView(e, app) }))
	r.POST("/admin/works/{id}", post(func(e *core.RequestEvent) error { return adminWorkEditSubmit(e, app) }))
	r.POST("/admin/works/{id}/delete", post(func(e *core.RequestEvent) error { return adminWorkDelete(e, app) }))
	r.POST("/admin/works/{id}/images/delete", post(func(e *core.RequestEvent) error { return adminWorkImageDelete(e, app) }))
	r.POST("/admin/works/{id}/images/move", post(func(e *core.RequestEvent) error { return adminWorkImageMove(e, app) }))

	// --- Admin: work links (HTMX) ---
	r.GET("/admin/works/{id}/links/new", pub(func(e *core.RequestEvent) error { return adminWorkLinkNewForm(e, app) }))
	r.POST("/admin/works/{id}/links", post(func(e *core.RequestEvent) error { return adminWorkLinkCreate(e, app) }))
	r.POST("/admin/works/{wid}/links/{lid}/delete", post(func(e *core.RequestEvent) error { return adminWorkLinkDelete(e, app) }))

	// --- Admin: press ---
	r.GET("/admin/press", pub(func(e *core.RequestEvent) error { return adminPressList(e, app) }))
	r.GET("/admin/press/new", pub(func(e *core.RequestEvent) error { return adminPressNewView(e, app) }))
	r.POST("/admin/press/new", post(func(e *core.RequestEvent) error { return adminPressNewSubmit(e, app) }))
	r.GET("/admin/press/{id}", pub(func(e *core.RequestEvent) error { return adminPressEditView(e, app) }))
	r.POST("/admin/press/{id}", post(func(e *core.RequestEvent) error { return adminPressEditSubmit(e, app) }))
	r.POST("/admin/press/{id}/delete", post(func(e *core.RequestEvent) error { return adminPressDelete(e, app) }))

	// --- Admin: settings ---
	r.GET("/admin/settings", pub(func(e *core.RequestEvent) error { return adminSettingsView(e, app) }))
	r.POST("/admin/settings", post(func(e *core.RequestEvent) error { return adminSettingsSubmit(e, app) }))
	r.POST("/admin/settings/hero_images/delete", post(func(e *core.RequestEvent) error { return adminHeroImageDelete(e, app) }))
	r.POST("/admin/settings/hero_images/move", post(func(e *core.RequestEvent) error { return adminHeroImageMove(e, app) }))
}
