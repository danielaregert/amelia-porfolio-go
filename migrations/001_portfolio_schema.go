package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Esquema del portfolio de Amelia Repetto:
//   site_settings  — configuración global (bio, contacto, redes)
//   sections       — categorías principales (Teatro, Performance, Cine, etc.)
//   works          — obras/proyectos individuales dentro de cada sección
//   work_links     — links asociados a una obra (trailer, dossier, prensa, etc.)
//   press          — notas de prensa independientes
func init() {
	m.Register(func(app core.App) error {

		// ---------- site_settings ----------
		// Singleton: siempre habrá un solo registro con la config del sitio.
		settings := core.NewBaseCollection("site_settings")
		settings.Fields.Add(&core.TextField{Name: "site_name", Required: true, Max: 200})
		settings.Fields.Add(&core.EditorField{Name: "bio_es"})
		settings.Fields.Add(&core.EditorField{Name: "bio_en"})
		settings.Fields.Add(&core.TextField{Name: "tagline", Max: 300})
		settings.Fields.Add(&core.EmailField{Name: "email"})
		settings.Fields.Add(&core.TextField{Name: "phone", Max: 40})
		settings.Fields.Add(&core.TextField{Name: "address", Max: 500})
		settings.Fields.Add(&core.URLField{Name: "instagram"})
		settings.Fields.Add(&core.URLField{Name: "facebook"})
		settings.Fields.Add(&core.URLField{Name: "youtube"})
		settings.Fields.Add(&core.URLField{Name: "vimeo"})
		settings.Fields.Add(&core.URLField{Name: "website"})
		settings.Fields.Add(&core.FileField{
			Name:      "profile_image",
			MaxSelect: 1,
			MaxSize:   10 * 1024 * 1024,
			MimeTypes: []string{"image/jpeg", "image/png", "image/webp"},
		})
		settings.Fields.Add(&core.FileField{
			Name:      "hero_image",
			MaxSelect: 1,
			MaxSize:   10 * 1024 * 1024,
			MimeTypes: []string{"image/jpeg", "image/png", "image/webp"},
		})
		settings.Fields.Add(&core.TextField{Name: "reel_url", Max: 500})
		settings.ListRule = strPtr("")
		settings.ViewRule = strPtr("")
		if err := app.Save(settings); err != nil {
			return err
		}

		// ---------- sections ----------
		// Las secciones principales: Performance, Teatro, Cine, Videoclips,
		// Talleres, Salud Mental
		sections := core.NewBaseCollection("sections")
		sections.Fields.Add(&core.TextField{Name: "name", Required: true, Max: 100})
		sections.Fields.Add(&core.TextField{Name: "slug", Required: true, Max: 100})
		sections.Fields.Add(&core.EditorField{Name: "description"})
		sections.Fields.Add(&core.NumberField{Name: "sort_order", OnlyInt: true})
		sections.Fields.Add(&core.BoolField{Name: "active"})
		sections.Fields.Add(&core.FileField{
			Name:      "cover_image",
			MaxSelect: 1,
			MaxSize:   10 * 1024 * 1024,
			MimeTypes: []string{"image/jpeg", "image/png", "image/webp"},
		})
		sections.AddIndex("idx_sections_slug", true, "slug", "")
		sections.AddIndex("idx_sections_order", false, "sort_order", "")
		sections.ListRule = strPtr("active = true")
		sections.ViewRule = strPtr("active = true")
		if err := app.Save(sections); err != nil {
			return err
		}

		// ---------- works ----------
		// Cada obra/proyecto: una entrada de teatro, performance, film, taller, etc.
		works := core.NewBaseCollection("works")
		works.Fields.Add(&core.RelationField{
			Name:          "section",
			Required:      true,
			MaxSelect:     1,
			CollectionId:  sections.Id,
			CascadeDelete: true,
		})
		works.Fields.Add(&core.TextField{Name: "title", Required: true, Max: 300})
		works.Fields.Add(&core.TextField{Name: "slug", Required: true, Max: 300})
		works.Fields.Add(&core.TextField{Name: "year", Max: 20})
		works.Fields.Add(&core.TextField{Name: "role", Max: 200})
		works.Fields.Add(&core.EditorField{Name: "description"})
		works.Fields.Add(&core.TextField{Name: "credits", Max: 2000})
		works.Fields.Add(&core.FileField{
			Name:      "images",
			MaxSelect: 12,
			MaxSize:   10 * 1024 * 1024,
			MimeTypes: []string{"image/jpeg", "image/png", "image/webp", "image/avif"},
		})
		works.Fields.Add(&core.NumberField{Name: "sort_order", OnlyInt: true})
		works.Fields.Add(&core.BoolField{Name: "active"})
		works.Fields.Add(&core.BoolField{Name: "featured"})
		works.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		works.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})
		works.AddIndex("idx_works_slug", true, "slug", "")
		works.AddIndex("idx_works_section", false, "section", "")
		works.AddIndex("idx_works_order", false, "sort_order", "")
		works.ListRule = strPtr("active = true && section.active = true")
		works.ViewRule = strPtr("active = true && section.active = true")
		if err := app.Save(works); err != nil {
			return err
		}

		// ---------- work_links ----------
		// Links asociados a una obra: trailer, obra completa, dossier, ficha, etc.
		workLinks := core.NewBaseCollection("work_links")
		workLinks.Fields.Add(&core.RelationField{
			Name:          "work",
			Required:      true,
			MaxSelect:     1,
			CollectionId:  works.Id,
			CascadeDelete: true,
		})
		workLinks.Fields.Add(&core.TextField{Name: "label", Required: true, Max: 200})
		workLinks.Fields.Add(&core.URLField{Name: "url", Required: true})
		workLinks.Fields.Add(&core.SelectField{
			Name:      "kind",
			Required:  true,
			MaxSelect: 1,
			Values:    []string{"trailer", "full_work", "dossier", "press", "external", "instagram"},
		})
		workLinks.Fields.Add(&core.NumberField{Name: "sort_order", OnlyInt: true})
		workLinks.AddIndex("idx_work_links_work", false, "work", "")
		workLinks.ListRule = strPtr("")
		workLinks.ViewRule = strPtr("")
		if err := app.Save(workLinks); err != nil {
			return err
		}

		// ---------- press ----------
		// Notas de prensa independientes (no asociadas a una obra específica).
		press := core.NewBaseCollection("press")
		press.Fields.Add(&core.TextField{Name: "title", Required: true, Max: 500})
		press.Fields.Add(&core.TextField{Name: "publication", Max: 200})
		press.Fields.Add(&core.URLField{Name: "url"})
		press.Fields.Add(&core.TextField{Name: "date", Max: 40})
		press.Fields.Add(&core.EditorField{Name: "excerpt"})
		press.Fields.Add(&core.NumberField{Name: "sort_order", OnlyInt: true})
		press.Fields.Add(&core.BoolField{Name: "active"})
		press.AddIndex("idx_press_order", false, "sort_order", "")
		press.ListRule = strPtr("active = true")
		press.ViewRule = strPtr("active = true")
		if err := app.Save(press); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		for _, name := range []string{
			"press", "work_links", "works", "sections", "site_settings",
		} {
			c, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				continue
			}
			if err := app.Delete(c); err != nil {
				return err
			}
		}
		return nil
	})
}

func strPtr(s string) *string  { return &s }
func floatPtr(f float64) *float64 { return &f }
