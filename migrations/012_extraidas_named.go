package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Imágenes nombradas por Amelia en extraidas_pdf/.
// Las obras sin imágenes las reciben; las que ya tienen, se añaden al final.
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

		attach := func(slug string, rels []string) error {
			w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug})
			if err != nil || w == nil {
				return nil
			}
			var files []*filesystem.File
			for _, rel := range rels {
				if f := mk(rel); f != nil {
					files = append(files, f)
				}
			}
			if len(files) == 0 {
				return nil
			}
			if len(w.GetStringSlice("images")) == 0 {
				w.Set("images", files)
			} else {
				w.Set("+images", files)
			}
			return app.Save(w)
		}

		// Obras sin imágenes → asignar
		toAttach := map[string][]string{
			"amor-mentiras-dinero":  {"extraidas_pdf/amor-mentias-dinero.jpg"},
			"ballux-instalacion":    {"extraidas_pdf/instalacion ballux.jpg"},
			"el-pendulo":            {"extraidas_pdf/el-pendulo.jpg"},
			"ensuenos":              {"extraidas_pdf/ensueños.jpg"},
			"la-debil-mental":       {"extraidas_pdf/la'debil mental.jpg"},
			"la-piramide":           {"extraidas_pdf/la piramide.jpg", "extraidas_pdf/lapiramida.jpg"},
			"picnic-2023":           {"extraidas_pdf/picnic 2022.jpg", "extraidas_pdf/picnic 2022 2.jpg", "extraidas_pdf/picnic 2022 3.jpg", "extraidas_pdf/picnic 2022 4.jpg"},
			"ezeiza-intervenciones": {"extraidas_pdf/salud-mental.jpg"},

			// Obras con imágenes → append
			"cronicas-ferreteras":    {"extraidas_pdf/cronicas ferreteras.jpg"},
			"la-liebre-y-la-tortuga": {"extraidas_pdf/laliebre y la tortuga.jpg"},
			"las-cuerdas":            {"extraidas_pdf/las-cuerdas.jpg"},
			"las-multitudes":         {"extraidas_pdf/las multitudes.jpg"},
			"todo-lo-que-veo-es-mio": {"extraidas_pdf/todo lo que veo es mio.jpg"},
		}

		for slug, rels := range toAttach {
			if err := attach(slug, rels); err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error { return nil })
}
