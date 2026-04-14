package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Carga de imágenes nuevas de fotos2026/ + fwdagregados/ (nombradas por slug),
// adjunta el video .mov y el dossier PDF a Romance de la Negra Rubia,
// y cambia la foto de perfil por foto-ame.jpeg.
func init() {
	m.Register(func(app core.App) error {
		root := os.Getenv("AMELIA_DESIGN_DIR")
		if root == "" {
			root = "/Users/danielaregert/Design/PORFOLIO AME"
		}

		// 1. Agregar campos video + dossier a works.
		works, err := app.FindCollectionByNameOrId("works")
		if err != nil {
			return err
		}
		if works.Fields.GetByName("video") == nil {
			works.Fields.Add(&core.FileField{
				Name:      "video",
				MaxSelect: 1,
				MaxSize:   200 * 1024 * 1024,
				MimeTypes: []string{"video/mp4", "video/webm", "video/quicktime", "video/x-quicktime"},
			})
		}
		if works.Fields.GetByName("dossier") == nil {
			works.Fields.Add(&core.FileField{
				Name:      "dossier",
				MaxSelect: 1,
				MaxSize:   50 * 1024 * 1024,
				MimeTypes: []string{"application/pdf"},
			})
		}
		if err := app.Save(works); err != nil {
			return err
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

		attachImages := func(slug string, rels []string) error {
			w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug})
			if err != nil || w == nil {
				return nil
			}
			var files []*filesystem.File
			for _, r := range rels {
				if f := mk(r); f != nil {
					files = append(files, f)
				}
			}
			if len(files) == 0 {
				return nil
			}
			w.Set("images", files)
			return app.Save(w)
		}

		// 2. Imágenes por obra (nombradas en fotos2026/).
		if err := attachImages("vitis-vinifera", []string{
			"fotos2026/vitis-vinifera.jpeg",
			"fotos2026/VITIS F3.jpg",
		}); err != nil {
			return err
		}
		if err := attachImages("bzd-hasta-los-huesos", []string{
			"fotos2026/bzd-hasta-los-huesos.JPG",
			"fotos2026/bzd-huesos.jpeg",
			"fotos2026/fwdagregados/bzd-hasta-los-huesos.jpeg",
			"fotos2026/fwdagregados/bzd-hasta-los-huesos-2.jpeg",
			"fotos2026/fwdagregados/bzd-hasta-los-huesos-3.jpeg",
		}); err != nil {
			return err
		}
		if err := attachImages("romance-negra-rubia", []string{
			"fotos2026/fwdagregados/romance-negra-rubia.jpeg",
			"fotos2026/fwdagregados/romance-negra-rubia-2.jpeg",
			"fotos2026/fwdagregados/romance-negra-rubia-flyer.jpeg",
		}); err != nil {
			return err
		}

		// 3. Video + dossier a Romance de la Negra Rubia.
		romance, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
			map[string]any{"s": "romance-negra-rubia"})
		if err == nil && romance != nil {
			if v := mk("fotos2026/fwdagregados/romance de la negra rubia video.mov"); v != nil {
				romance.Set("video", v)
			}
			if d := mk("fotos2026/fwdagregados/Romance de la negra rubia DOSSIER.pdf"); d != nil {
				romance.Set("dossier", d)
			}
			if err := app.Save(romance); err != nil {
				return err
			}
		}

		// 4. Foto de perfil → foto-ame.jpeg.
		settings, err := app.FindFirstRecordByFilter("site_settings", "id != ''")
		if err == nil && settings != nil {
			if f := mk("fotos2026/fwdagregados/foto-ame.jpeg"); f != nil {
				settings.Set("profile_image", f)
				if f2 := mk("fotos2026/fwdagregados/foto-ame.jpeg"); f2 != nil {
					settings.Set("hero_image", f2)
				}
				if err := app.Save(settings); err != nil {
					return err
				}
			}
		}

		return nil
	}, func(app core.App) error { return nil })
}
