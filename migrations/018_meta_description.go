package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Agrega meta_description + meta_description_en a site_settings para SEO / OG.
func init() {
	m.Register(func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}
		if col.Fields.GetByName("meta_description") == nil {
			col.Fields.Add(&core.TextField{Name: "meta_description", Max: 300})
		}
		if col.Fields.GetByName("meta_description_en") == nil {
			col.Fields.Add(&core.TextField{Name: "meta_description_en", Max: 300})
		}
		if err := app.Save(col); err != nil {
			return err
		}

		// Seed inicial si están vacíos.
		if rec, err := app.FindFirstRecordByFilter("site_settings", "id != ''"); err == nil && rec != nil {
			if rec.GetString("meta_description") == "" {
				rec.Set("meta_description",
					"Amelia Repetto — actriz, directora de teatro y psiquiatra especializada en arte y salud mental. Residente en Madrid. Porfolio de obras de teatro, performance, cine, videoclips, talleres y proyectos de salud mental.")
			}
			if rec.GetString("meta_description_en") == "" {
				rec.Set("meta_description_en",
					"Amelia Repetto — actress, theatre director, and psychiatrist specialized in art and mental health. Based in Madrid. Portfolio of theatre, performance, film, music videos, workshops, and mental health projects.")
			}
			_ = app.Save(rec)
		}
		return nil
	}, func(app core.App) error { return nil })
}
