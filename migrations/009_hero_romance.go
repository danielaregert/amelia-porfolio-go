package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Reemplaza la 3ra imagen del hero por una de Romance de la Negra Rubia.
func init() {
	m.Register(func(app core.App) error {
		root := os.Getenv("AMELIA_DESIGN_DIR")
		if root == "" {
			root = "/Users/danielaregert/Design/PORFOLIO AME"
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
		var files []*filesystem.File
		for _, rel := range []string{
			"fotos2026/fwdagregados/foto-ame.jpeg",
			"fotos2026/image00012.jpeg",
			"fotos2026/fwdagregados/romance-negra-rubia.jpeg",
		} {
			if f := mk(rel); f != nil {
				files = append(files, f)
			}
		}
		if len(files) == 0 {
			return nil
		}
		rec.Set("hero_images", files)
		return app.Save(rec)
	}, func(app core.App) error { return nil })
}
