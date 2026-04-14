package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Actualiza links de obras según URLs reales extraídas del .indd
// y correcciones aportadas por la artista.
func init() {
	m.Register(func(app core.App) error {
		workLinks, err := app.FindCollectionByNameOrId("work_links")
		if err != nil {
			return err
		}

		// Borra todos los links de una obra dada por slug.
		purge := func(slug string) error {
			w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug})
			if err != nil || w == nil {
				return nil
			}
			links, err := app.FindRecordsByFilter("work_links",
				"work = {:w}", "", 500, 0, map[string]any{"w": w.Id})
			if err != nil {
				return err
			}
			for _, l := range links {
				if err := app.Delete(l); err != nil {
					return err
				}
			}
			return nil
		}

		addLinks := func(slug string, links []linkSeed) error {
			w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug})
			if err != nil || w == nil {
				return nil
			}
			for i, l := range links {
				lr := core.NewRecord(workLinks)
				lr.Set("work", w.Id)
				lr.Set("label", l.label)
				lr.Set("url", l.url)
				lr.Set("kind", l.kind)
				lr.Set("sort_order", i+1)
				if err := app.Save(lr); err != nil {
					return err
				}
			}
			return nil
		}

		// ---- Videoclips (URLs reales del .indd) ----
		videoclips := map[string]string{
			"no-todo-es-color-de-rosa": "https://www.youtube.com/watch?v=uFMl4ZCx_tc",
			"casa-roja":                "https://www.youtube.com/watch?v=Md1GvylcXF0&list=RDMd1GvylcXF0&start_radio=1",
			"sensaciones":              "https://www.youtube.com/watch?v=RrHY-xE2Yos&list=RDRrHY-xE2Yos&start_radio=1",
			"linda":                    "https://www.youtube.com/watch?v=uiy09JNsZ4A",
		}
		for slug, url := range videoclips {
			if err := addLinks(slug, []linkSeed{
				{"Ver videoclip", url, "full_work"},
			}); err != nil {
				return err
			}
		}

		// ---- Romance de la Negra Rubia: reemplazar placeholders ----
		if err := purge("romance-negra-rubia"); err != nil {
			return err
		}
		if err := addLinks("romance-negra-rubia", []linkSeed{
			{"Trailer", "https://youtu.be/c2RhG649Lxk?is=akObc0JxCAL4gFkZ", "trailer"},
			{"Obra completa", "https://youtu.be/5eelF51Q_44?is=2O0UaixCt5MehCc5", "full_work"},
		}); err != nil {
			return err
		}

		// ---- BZD hasta los huesos: reemplazar placeholder ----
		if err := purge("bzd-hasta-los-huesos"); err != nil {
			return err
		}
		if err := addLinks("bzd-hasta-los-huesos", []linkSeed{
			{"Link a investigación", "https://youtu.be/twBZih2Lf1k?is=Z1Nt5kaT4ZfnGu-J", "full_work"},
		}); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		// Down: no-op (no perdemos los videoclips si rollback).
		return nil
	})
}
