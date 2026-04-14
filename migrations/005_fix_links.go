package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Limpia placeholders (vimeo.com/..., drive.google.com/..., otro.com/..., etc.)
// y reemplaza por los links reales extraídos del archivo .indd del porfolio.
// También elimina duplicados.
func init() {
	m.Register(func(app core.App) error {

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

		workLinks, err := app.FindCollectionByNameOrId("work_links")
		if err != nil {
			return err
		}
		set := func(slug string, links []linkSeed) error {
			if err := purge(slug); err != nil {
				return err
			}
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

		// ============ TEATRO ============
		if err := set("romance-negra-rubia", []linkSeed{
			{"Trailer", "https://youtu.be/c2RhG649Lxk?is=akObc0JxCAL4gFkZ", "trailer"},
			{"Obra completa", "https://youtu.be/5eelF51Q_44?is=2O0UaixCt5MehCc5", "full_work"},
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra57172-el-romance-de-la-negra-rubia", "external"},
		}); err != nil {
			return err
		}

		if err := set("la-debil-mental", []linkSeed{
			{"Suplemento Cultural Ñ Clarín", "https://www.clarin.com/revista-enie/literatura/ariana-harwicz-racismo-pedofilia-maternidad-salvaje_0_ia38yau4y.html", "press"},
			{"Trabajos previos sobre la novela (Radiocut)", "https://ar.radiocut.fm/audiocut/lectura-cristina-banegas-debil-mental-ariana-harwicz/", "external"},
			{"Instagram @ladebilmentalteatro", "https://www.instagram.com/ladebilmentalteatro/", "instagram"},
		}); err != nil {
			return err
		}

		if err := set("las-cuerdas", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra68053-las-cuerdas", "external"},
			{"Nota Página 12", "https://www.pagina12.com.ar/327284-las-cuerdas-de-ana-schimelman", "press"},
			{"Nota Fervor", "https://fervor.com.ar/los-bordes-vinculares/", "press"},
		}); err != nil {
			return err
		}

		if err := set("el-pendulo", []linkSeed{
			{"Festival Fauna", "https://fauna.una.edu.ar/obra/el-pendulo_234_6", "external"},
		}); err != nil {
			return err
		}

		if err := set("amor-mentiras-dinero", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra59713-amor-mentiras-dinero", "external"},
		}); err != nil {
			return err
		}

		if err := set("ensuenos", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra59132-en-suenos-sala-1", "external"},
		}); err != nil {
			return err
		}

		if err := set("el-tilo", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra56886-el-tilo", "external"},
		}); err != nil {
			return err
		}

		if err := set("la-liebre-y-la-tortuga", []linkSeed{
			{"Ficha (Teatro Cervantes)", "https://www.teatrocervantes.gob.ar/obra/la-liebre-y-la-tortuga/", "external"},
		}); err != nil {
			return err
		}

		if err := set("zoom-in-90s", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra39085-zoom-in-90s", "external"},
		}); err != nil {
			return err
		}

		if err := set("las-multitudes", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra24764-las-multitudes", "external"},
			{"Trailer (steirischerherbst)", "https://www.youtube.com/watch?v=gQgEcAQPJtQ&ab_channel=steirischerherbst", "trailer"},
			{"Nota La Nación", "https://www.lanacion.com.ar/espectaculos/teatro/las-multitudes-nid1493327/", "press"},
			{"Nota Página 12 (Radar)", "https://www.pagina12.com.ar/diario/suplementos/radar/9-8075-2012-07-15.html", "press"},
			{"Farsa mag", "https://farsamag.com.ar/obras/las-multitudes/", "press"},
			{"Revista Otra Parte", "https://www.revistaotraparte.com/seccion/ensayo-teoria/", "press"},
		}); err != nil {
			return err
		}

		if err := set("cinefilia", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra32425-cinefilia", "external"},
		}); err != nil {
			return err
		}

		if err := set("vispera-de-elecciones", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra15766-vispera-de-elecciones", "external"},
			{"Nota La Nación", "https://www.lanacion.com.ar/espectaculos/teatro/vispera-de-elecciones-nid1422453/", "press"},
		}); err != nil {
			return err
		}

		if err := set("juicio-lady-macbeth", []linkSeed{
			{"Ficha (alternativa teatral)", "http://www.alternativateatral.com/obra15564-el-juicio-de-lady-macbeth", "external"},
			{"Blog UMDH", "http://umdh.blogspot.com/2010/04/el-juicio-de-lady-macbeth-llega-san.html?m=1", "press"},
		}); err != nil {
			return err
		}

		if err := set("lo-que-se-dice", []linkSeed{
			{"Video / Ficha (alternativa teatral)", "http://www.alternativateatral.com/video9643-lo-que-se-dice", "external"},
		}); err != nil {
			return err
		}

		// Vitis Vinífera (2025): sin URLs reales en .indd — dejar vacío, cargar desde admin.
		if err := purge("vitis-vinifera"); err != nil {
			return err
		}

		// La Pirámide, Instalación Ballux: sin URLs en el PDF. Vaciar si hay placeholders.
		for _, s := range []string{"la-piramide", "ballux-instalacion"} {
			if err := purge(s); err != nil {
				return err
			}
		}

		// ============ CINE ============
		if err := set("inconsciente-colectivo", []linkSeed{
			{"Ficha Prime Video", "https://www.primevideo.com/detail/Inconsciente-Colectivo/0MW848O4B7T2FWGJJF5HYKRL27", "external"},
		}); err != nil {
			return err
		}

		if err := set("noemi-gold", []linkSeed{
			{"Imdb", "https://www.imdb.com/title/tt8548708/", "external"},
			{"Filmaffinity", "https://www.filmaffinity.com/es/film396367.html", "external"},
			{"Mubi", "https://mubi.com/es/films/noemi-gold/cast", "external"},
			{"Topic", "https://www.topic.com/indie-filmmaking-is-not-dead-a-conversation-with-noemi-gold-director-dan-rubenstein", "press"},
			{"Deadline", "https://deadline.com/2020/11/topic-dan-rubensteins-argentinian-drama-noemi-gold-1234612997/", "press"},
			{"Cineramaplus", "https://cineramaplus.com.ar/critica-noemi-gold-2019-de-dan-rubenstein-bafici/", "press"},
			{"Rogers Movie Nation", "https://rogersmovienation.com/2020/11/21/movie-review-argentine-and-in-need-of-an-abortion-noemi-gold/", "press"},
			{"Vice", "https://www.vice.com/es/article/3kg3g5/no-tener-la-libertad-para-poder-decidir-sobre-tu-cuerpo-me-parece-ridiculo-sumamente-peligroso-y-clasista", "press"},
			{"Cinéfiloserial", "https://cinefiloserial.com.ar/21-bafici-noemi-gold-de-dan-rubenstein-2019/", "press"},
		}); err != nil {
			return err
		}

		if err := set("cronicas-ferreteras", []linkSeed{
			{"Cinear (play.cine.ar)", "https://play.cine.ar/INCAA/produccion/", "external"},
		}); err != nil {
			return err
		}

		// Todo lo que veo es mío, El éxtasis, Veredas: sin URL clara en .indd — vaciar.
		for _, s := range []string{"todo-lo-que-veo-es-mio", "extasis-santa-teresa", "veredas"} {
			if err := purge(s); err != nil {
				return err
			}
		}

		// ============ VIDEOCLIPS ============
		if err := set("no-todo-es-color-de-rosa", []linkSeed{
			{"Ver videoclip", "https://www.youtube.com/watch?v=uFMl4ZCx_tc", "full_work"},
		}); err != nil {
			return err
		}
		if err := set("casa-roja", []linkSeed{
			{"Ver videoclip", "https://www.youtube.com/watch?v=Md1GvylcXF0&list=RDMd1GvylcXF0&start_radio=1", "full_work"},
		}); err != nil {
			return err
		}
		if err := set("sensaciones", []linkSeed{
			{"Ver videoclip", "https://www.youtube.com/watch?v=RrHY-xE2Yos&list=RDRrHY-xE2Yos&start_radio=1", "full_work"},
		}); err != nil {
			return err
		}
		if err := set("linda", []linkSeed{
			{"Ver videoclip", "https://www.youtube.com/watch?v=uiy09JNsZ4A", "full_work"},
		}); err != nil {
			return err
		}

		// ============ PERFORMANCE ============
		if err := set("bzd-hasta-los-huesos", []linkSeed{
			{"Link a investigación", "https://youtu.be/twBZih2Lf1k?is=Z1Nt5kaT4ZfnGu-J", "full_work"},
		}); err != nil {
			return err
		}

		if err := set("castle-crossed-destinies", []linkSeed{
			{"Ficha Museo Reina Sofía", "https://www.museoreinasofia.es/exposiciones/leonor-serrano-rivas", "external"},
		}); err != nil {
			return err
		}

		if err := set("picnic-2023", []linkSeed{
			{"Ficha Museo Reina Sofía (PICNIC 2022)", "https://www.museoreinasofia.es/actividades/picnic-barrio-2022", "external"},
		}); err != nil {
			return err
		}

		// Obras sin URL real confirmada en el .indd — dejamos sin links.
		for _, s := range []string{
			"bzd-el-entierro",
			"artificio-psicosis",
			"artificio-cuerpos-dociles",
			"stop-exclusion-sanitaria",
		} {
			if err := purge(s); err != nil {
				return err
			}
		}

		// ============ SALUD MENTAL ============
		// Sin URLs reales identificadas en .indd — vaciar placeholders.
		for _, s := range []string{"ezeiza-intervenciones", "el-juicio-f22"} {
			if err := purge(s); err != nil {
				return err
			}
		}

		// ============ site_settings: redes ============
		settings, err := app.FindFirstRecordByFilter("site_settings", "id != ''")
		if err == nil && settings != nil {
			settings.Set("instagram", "https://www.instagram.com/amelia.repetto/?hl=es")
			settings.Set("facebook", "https://www.facebook.com/amelia.repetto/")
			_ = app.Save(settings)
		}

		// ============ press ============
		press, err := app.FindRecordsByFilter("press", "id != ''", "+sort_order", 100, 0)
		if err == nil {
			for _, p := range press {
				switch p.GetString("title") {
				case "De Buenos Aires a Madrid sin escalas: la psiquiatra argentina que abandonó una profesión de 11 años por amor al teatro":
					p.Set("url", "https://www.elmundo.es/madrid/2023/09/04/64f20d56fc6c834e328b457c.html")
					_ = app.Save(p)
				case "Amelia Repetto presentó “Artificio para atravesar la psicosis y los cuerpos dóciles” en la Asociación ATLAS":
					p.Set("url", "https://www.masescena.es/index.php/noticias/otras-disciplinas/8236-amelia-repetto-presento-artificio-para-atravesar-la-psicosis-y-los-cuerpos-dociles-en-la-asociacion-atlas-con-participacion-ciudadana-y-dentro-del-festival-tara")
					_ = app.Save(p)
				}
			}
		}

		return nil
	}, func(app core.App) error { return nil })
}
