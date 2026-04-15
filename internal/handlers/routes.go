package handlers

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func RegisterRoutes(r *router.Router[*core.RequestEvent], app *pocketbase.PocketBase) {
	// --- Public (Phase 2) ---
	r.GET("/", func(e *core.RequestEvent) error {
		return publicHome(e, app)
	})
	r.GET("/en", func(e *core.RequestEvent) error {
		return publicHome(e, app)
	})
	r.GET("/en/", func(e *core.RequestEvent) error {
		return publicHome(e, app)
	})

	// --- Admin: auth ---
	r.GET("/admin/login", func(e *core.RequestEvent) error {
		return adminLoginView(e, app)
	})
	r.POST("/admin/login", func(e *core.RequestEvent) error {
		return adminLoginSubmit(e, app)
	})
	r.POST("/admin/logout", func(e *core.RequestEvent) error {
		return adminLogout(e, app)
	})

	// --- Admin: sections ---
	r.GET("/admin", func(e *core.RequestEvent) error {
		return adminSectionsList(e, app)
	})
	r.GET("/admin/sections/new", func(e *core.RequestEvent) error {
		return adminSectionNewView(e, app)
	})
	r.POST("/admin/sections/new", func(e *core.RequestEvent) error {
		return adminSectionNewSubmit(e, app)
	})
	r.GET("/admin/sections/{id}", func(e *core.RequestEvent) error {
		return adminSectionEditView(e, app)
	})
	r.POST("/admin/sections/{id}", func(e *core.RequestEvent) error {
		return adminSectionEditSubmit(e, app)
	})
	r.POST("/admin/sections/{id}/toggle", func(e *core.RequestEvent) error {
		return adminSectionToggle(e, app)
	})
	r.POST("/admin/sections/{id}/delete", func(e *core.RequestEvent) error {
		return adminSectionDelete(e, app)
	})

	// --- Admin: works ---
	r.GET("/admin/works", func(e *core.RequestEvent) error {
		return adminWorksList(e, app)
	})
	r.GET("/admin/works/new", func(e *core.RequestEvent) error {
		return adminWorkNewView(e, app)
	})
	r.POST("/admin/works/new", func(e *core.RequestEvent) error {
		return adminWorkNewSubmit(e, app)
	})
	r.GET("/admin/works/{id}", func(e *core.RequestEvent) error {
		return adminWorkEditView(e, app)
	})
	r.POST("/admin/works/{id}", func(e *core.RequestEvent) error {
		return adminWorkEditSubmit(e, app)
	})
	r.POST("/admin/works/{id}/delete", func(e *core.RequestEvent) error {
		return adminWorkDelete(e, app)
	})
	r.POST("/admin/works/{id}/images/delete", func(e *core.RequestEvent) error {
		return adminWorkImageDelete(e, app)
	})
	r.POST("/admin/works/{id}/images/move", func(e *core.RequestEvent) error {
		return adminWorkImageMove(e, app)
	})

	// --- Admin: work links (HTMX) ---
	r.GET("/admin/works/{id}/links/new", func(e *core.RequestEvent) error {
		return adminWorkLinkNewForm(e, app)
	})
	r.POST("/admin/works/{id}/links", func(e *core.RequestEvent) error {
		return adminWorkLinkCreate(e, app)
	})
	r.POST("/admin/works/{wid}/links/{lid}/delete", func(e *core.RequestEvent) error {
		return adminWorkLinkDelete(e, app)
	})

	// --- Admin: press ---
	r.GET("/admin/press", func(e *core.RequestEvent) error {
		return adminPressList(e, app)
	})
	r.GET("/admin/press/new", func(e *core.RequestEvent) error {
		return adminPressNewView(e, app)
	})
	r.POST("/admin/press/new", func(e *core.RequestEvent) error {
		return adminPressNewSubmit(e, app)
	})
	r.GET("/admin/press/{id}", func(e *core.RequestEvent) error {
		return adminPressEditView(e, app)
	})
	r.POST("/admin/press/{id}", func(e *core.RequestEvent) error {
		return adminPressEditSubmit(e, app)
	})
	r.POST("/admin/press/{id}/delete", func(e *core.RequestEvent) error {
		return adminPressDelete(e, app)
	})

	// --- Admin: settings ---
	r.GET("/admin/settings", func(e *core.RequestEvent) error {
		return adminSettingsView(e, app)
	})
	r.POST("/admin/settings", func(e *core.RequestEvent) error {
		return adminSettingsSubmit(e, app)
	})
	r.POST("/admin/settings/hero_images/delete", func(e *core.RequestEvent) error {
		return adminHeroImageDelete(e, app)
	})
	r.POST("/admin/settings/hero_images/move", func(e *core.RequestEvent) error {
		return adminHeroImageMove(e, app)
	})
}
