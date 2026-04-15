package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Agrega og_image dedicado para social preview (fallback a primera del hero).
func init() {
	m.Register(func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}
		if col.Fields.GetByName("og_image") == nil {
			col.Fields.Add(&core.FileField{
				Name:      "og_image",
				MaxSelect: 1,
				MaxSize:   8 * 1024 * 1024,
				MimeTypes: []string{"image/jpeg", "image/png", "image/webp"},
			})
			return app.Save(col)
		}
		return nil
	}, func(app core.App) error { return nil })
}
