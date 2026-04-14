package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Expande hero_images a 8 slots y agrega vitis cartel + bzd-2.
func init() {
	m.Register(func(app core.App) error {
		root := os.Getenv("AMELIA_DESIGN_DIR")
		if root == "" {
			root = "/Users/danielaregert/Design/PORFOLIO AME"
		}

		// Subir max slots a 8.
		col, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}
		if fld := col.Fields.GetByName("hero_images"); fld != nil {
			if f, ok := fld.(*core.FileField); ok {
				f.MaxSelect = 8
				if err := app.Save(col); err != nil {
					return err
				}
			}
		}

		mk := func(rel string) *filesystem.File {
			full := filepath.Join(root, rel)
			if _, err := os.Stat(full); err != nil {
				return nil
			}
			f, err := filesystem.NewFileFromPath(full)
			if err != nil {
				return nil
			}
			return f
		}

		rec, err := app.FindFirstRecordByFilter("site_settings", "id != ''")
		if err != nil || rec == nil {
			return nil
		}
		var extra []*filesystem.File
		for _, rel := range []string{
			"fotos2026/VITIS F3.jpg",
			"fotos2026/fwdagregados/bzd-hasta-los-huesos-2.jpeg",
		} {
			if f := mk(rel); f != nil {
				extra = append(extra, f)
			}
		}
		if len(extra) == 0 {
			return nil
		}
		rec.Set("+hero_images", extra)
		return app.Save(rec)
	}, func(app core.App) error { return nil })
}
