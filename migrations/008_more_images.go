package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Más imágenes extraídas de los .indd 2021 y 2023 (obras más viejas).
// Mapeo confirmado por nombre de archivo → slug de obra.
func init() {
	m.Register(func(app core.App) error {
		root := os.Getenv("AMELIA_DESIGN_DIR")
		if root == "" {
			root = "/Users/danielaregert/Design/PORFOLIO AME"
		}

		attach := func(slug string, rels []string, mode string) error {
			w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug})
			if err != nil || w == nil {
				return nil
			}
			var files []*filesystem.File
			for _, rel := range rels {
				full := filepath.Join(root, rel)
				if _, err := os.Stat(full); err != nil {
					continue
				}
				f, err := filesystem.NewFileFromPath(full)
				if err != nil {
					continue
				}
				files = append(files, f)
			}
			if len(files) == 0 {
				return nil
			}
			if mode == "append" {
				w.Set("+images", files)
			} else {
				w.Set("images", files)
			}
			return app.Save(w)
		}

		// Reemplazar: obras sin imágenes
		if err := attach("cinefilia", []string{
			"FOTOS/CINEFILIA.jpg",
			"FOTOS/CINEFILIA2.jpg",
		}, "replace"); err != nil {
			return err
		}
		if err := attach("la-liebre-y-la-tortuga", []string{
			"FOTOS/la liebre y la tortuga.jpg",
		}, "replace"); err != nil {
			return err
		}
		if err := attach("lo-que-se-dice", []string{
			"FOTOS/lo que se dice.jpg",
		}, "replace"); err != nil {
			return err
		}
		if err := attach("vispera-de-elecciones", []string{
			"FOTOS/vispera de elecciones.jpg",
			"FOTOS/vispera.jpg",
		}, "replace"); err != nil {
			return err
		}
		if err := attach("juicio-lady-macbeth", []string{
			"FOTOS/el juicio.jpg",
		}, "replace"); err != nil {
			return err
		}
		if err := attach("las-cuerdas", []string{
			"las cuerdas 2.jpg",
		}, "replace"); err != nil {
			return err
		}
		if err := attach("el-juicio-f22", []string{
			"FOTOS/corto-ame.JPG",
		}, "replace"); err != nil {
			return err
		}

		// Agregar: las-multitudes ya tiene 1, sumar 2 más
		if err := attach("las-multitudes", []string{
			"FOTOS/las multitudes 2.jpg",
			"FOTOS/las multitudes 3.jpg",
		}, "append"); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error { return nil })
}
