package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Adjunta las imágenes del .indd (links) a cada obra.
// Mapeo derivado de los nombres de archivo/carpetas extraídos del
// PORFOLIO-2026-indesign.indd — ver documentación en README.
func init() {
	m.Register(func(app core.App) error {
		root := os.Getenv("AMELIA_DESIGN_DIR")
		if root == "" {
			root = "/Users/danielaregert/Design/PORFOLIO AME"
		}

		mapping := map[string][]string{
			// Teatro
			"vitis-vinifera":       {"fotos2026/VITIS F3.jpg"},
			"las-multitudes":       {"FOTOS/las multitudes.jpg"},
			// Cine
			"inconsciente-colectivo": {"FOTOS/INCONCIENTE COLECTIVO.jpg"},
			"noemi-gold":             {"noemi gold.jpg"},
			"cronicas-ferreteras":    {"FOTOS/cronicas ferreteras.jpg"},
			"todo-lo-que-veo-es-mio": {"FOTOS/todo l q veo es mio.jpg"},
			"extasis-santa-teresa":   {"extasis.jpg"},
			// Performance
			"castle-crossed-destinies": {
				"FOTOS/castle/image_6483441 (2).JPG",
				"FOTOS/castle/image_6483441 (4).JPG",
				"FOTOS/castle/image_6483441 (5).JPG",
				"FOTOS/castle/image_6483441 (6).JPG",
				"FOTOS/castle/image_6483441 (7).JPG",
			},
			"artificio-psicosis": {
				"FOTOS/artificio/1.jpg",
				"FOTOS/artificio/2.jpg",
				"FOTOS/artificio/5.jpg",
			},
			"artificio-cuerpos-dociles": {
				"FOTOS/artificio/1.jpg",
				"FOTOS/artificio/2.jpg",
			},
			"stop-exclusion-sanitaria": {
				"FOTOS/stop/1.jpg",
				"FOTOS/stop/2.jpg",
				"FOTOS/stop/3.jpg",
				"FOTOS/stop/4.jpg",
				"FOTOS/stop/6.jpg",
			},
			// Talleres
			"taller-la-quimera": {
				"FOTOS/talleres/flyer.jpg",
				"FOTOS/talleres/image_6483441 (2).JPG",
				"FOTOS/talleres/image_6483441 (3).JPG",
			},
		}

		for slug, rels := range mapping {
			work, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug})
			if err != nil || work == nil {
				continue
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
				continue
			}
			work.Set("images", files)
			if err := app.Save(work); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		// Down: vaciar images de las obras afectadas.
		slugs := []string{
			"vitis-vinifera", "las-multitudes", "inconsciente-colectivo",
			"noemi-gold", "cronicas-ferreteras", "todo-lo-que-veo-es-mio",
			"extasis-santa-teresa", "castle-crossed-destinies",
			"artificio-psicosis", "artificio-cuerpos-dociles",
			"stop-exclusion-sanitaria", "taller-la-quimera",
		}
		for _, s := range slugs {
			w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": s})
			if err != nil || w == nil {
				continue
			}
			w.Set("images", nil)
			_ = app.Save(w)
		}
		return nil
	})
}
