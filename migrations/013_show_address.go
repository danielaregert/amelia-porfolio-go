package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Agrega show_address (bool) a site_settings — oculta la dirección pública.
func init() {
	m.Register(func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}
		if col.Fields.GetByName("show_address") == nil {
			col.Fields.Add(&core.BoolField{Name: "show_address"})
			if err := app.Save(col); err != nil {
				return err
			}
		}
		rec, err := app.FindFirstRecordByFilter("site_settings", "id != ''")
		if err == nil && rec != nil {
			rec.Set("show_address", true)
			_ = app.Save(rec)
		}
		return nil
	}, func(app core.App) error { return nil })
}
