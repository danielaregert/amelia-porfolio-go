package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Agrega hero_images (multi-file) para el carrusel del hero.
func init() {
	m.Register(func(app core.App) error {
		root := os.Getenv("AMELIA_DESIGN_DIR")
		if root == "" {
			root = "/Users/danielaregert/Design/PORFOLIO AME"
		}

		settingsCol, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}
		if settingsCol.Fields.GetByName("hero_images") == nil {
			settingsCol.Fields.Add(&core.FileField{
				Name:      "hero_images",
				MaxSelect: 3,
				MaxSize:   15 * 1024 * 1024,
				MimeTypes: []string{"image/jpeg", "image/png", "image/webp"},
			})
			if err := app.Save(settingsCol); err != nil {
				return err
			}
		}

		rec, err := app.FindFirstRecordByFilter("site_settings", "id != ''")
		if err != nil || rec == nil {
			return nil
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

		var files []*filesystem.File
		for _, rel := range []string{
			"fotos2026/fwdagregados/foto-ame.jpeg",
			"fotos2026/image00012.jpeg",
			"extraidas_pdf/img-000.jpg",
		} {
			if f := mk(rel); f != nil {
				files = append(files, f)
			}
		}
		if len(files) > 0 {
			rec.Set("hero_images", files)
			if err := app.Save(rec); err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error { return nil })
}
