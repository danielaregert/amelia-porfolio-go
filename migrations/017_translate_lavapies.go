package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Traducciones EN de obras cargadas desde el admin (no seed).
// Solo setea los campos _en vacíos — respeta cualquier traducción previa.
func init() {
	m.Register(func(app core.App) error {
		translate := func(slug string, translations map[string]string) error {
			w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug})
			if err != nil || w == nil {
				return nil
			}
			changed := false
			for field, val := range translations {
				if w.GetString(field) == "" && val != "" {
					w.Set(field, val)
					changed = true
				}
			}
			if !changed {
				return nil
			}
			return app.Save(w)
		}

		if err := translate("lavapies", map[string]string{
			"title_en":       "Lavapiés by Fernando Ferrer",
			"role_en":        "Lead actress",
			"description_en": "TEATRO DEL BARRIO",
		}); err != nil {
			return err
		}

		if err := translate("teatro-espeluznante", map[string]string{
			"title_en":       "Presentation of the Manifesto of Spine-Chilling Theatre by Francesca Giordano",
			"role_en":        "performer",
			"description_en": "Performance in Bariloche, Patagonia, Argentina.",
			"credits_en":     "Libros Drama",
		}); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error { return nil })
}
