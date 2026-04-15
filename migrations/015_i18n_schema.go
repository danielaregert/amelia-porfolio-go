package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Agrega campos _en para soporte bilingüe ES/EN.
// bio_en ya existe en site_settings (definido en la migración 001).
func init() {
	m.Register(func(app core.App) error {
		addText := func(colName, field string, max int) error {
			col, err := app.FindCollectionByNameOrId(colName)
			if err != nil {
				return err
			}
			if col.Fields.GetByName(field) != nil {
				return nil
			}
			col.Fields.Add(&core.TextField{Name: field, Max: max})
			return app.Save(col)
		}
		addEditor := func(colName, field string) error {
			col, err := app.FindCollectionByNameOrId(colName)
			if err != nil {
				return err
			}
			if col.Fields.GetByName(field) != nil {
				return nil
			}
			col.Fields.Add(&core.EditorField{Name: field})
			return app.Save(col)
		}

		// site_settings
		if err := addText("site_settings", "tagline_en", 300); err != nil {
			return err
		}

		// sections
		if err := addText("sections", "name_en", 100); err != nil {
			return err
		}
		if err := addEditor("sections", "description_en"); err != nil {
			return err
		}

		// works
		if err := addText("works", "title_en", 300); err != nil {
			return err
		}
		if err := addText("works", "role_en", 200); err != nil {
			return err
		}
		if err := addEditor("works", "description_en"); err != nil {
			return err
		}
		if err := addText("works", "credits_en", 2000); err != nil {
			return err
		}

		// press
		if err := addText("press", "title_en", 500); err != nil {
			return err
		}
		if err := addEditor("press", "excerpt_en"); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error { return nil })
}
