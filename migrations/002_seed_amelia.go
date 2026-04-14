package migrations

import (
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Semilla con el contenido del porfolio de Amelia Repetto 2026.
// Textos tomados del PDF PORFOLIO-2026.pdf.
// Imagen de perfil/hero: extraidas_pdf/img-000.jpg de la carpeta Design.
func init() {
	m.Register(func(app core.App) error {

		designDir := os.Getenv("AMELIA_DESIGN_DIR")
		if designDir == "" {
			designDir = "/Users/danielaregert/Design/PORFOLIO AME"
		}

		// ---------- site_settings ----------
		settingsCol, err := app.FindCollectionByNameOrId("site_settings")
		if err != nil {
			return err
		}
		settings := core.NewRecord(settingsCol)
		settings.Set("site_name", "Amelia Repetto")
		settings.Set("tagline", "porfolio artístico")
		settings.Set("bio_es", bioES)
		settings.Set("bio_en", bioEN)
		settings.Set("email", "amelia.repetto@gmail.com")
		settings.Set("phone", "666 137 644")
		settings.Set("address", "Av Marqués de Corbera 62. Piso 3E CP 28017, Madrid")
		settings.Set("instagram", "https://instagram.com/amelia.repetto")
		settings.Set("facebook", "https://facebook.com/amelia.repetto")

		coverPath := filepath.Join(designDir, "extraidas_pdf", "img-000.jpg")
		if _, err := os.Stat(coverPath); err == nil {
			if f, err := filesystem.NewFileFromPath(coverPath); err == nil {
				settings.Set("profile_image", f)
			}
			if f, err := filesystem.NewFileFromPath(coverPath); err == nil {
				settings.Set("hero_image", f)
			}
		}
		if err := app.Save(settings); err != nil {
			return err
		}

		// ---------- sections ----------
		sectionsCol, err := app.FindCollectionByNameOrId("sections")
		if err != nil {
			return err
		}
		sectionIDs := map[string]string{}
		for i, s := range sectionsSeed {
			r := core.NewRecord(sectionsCol)
			r.Set("name", s.name)
			r.Set("slug", s.slug)
			r.Set("description", s.desc)
			r.Set("sort_order", i+1)
			r.Set("active", true)
			if err := app.Save(r); err != nil {
				return err
			}
			sectionIDs[s.slug] = r.Id
		}

		// ---------- works ----------
		worksCol, err := app.FindCollectionByNameOrId("works")
		if err != nil {
			return err
		}
		workLinksCol, err := app.FindCollectionByNameOrId("work_links")
		if err != nil {
			return err
		}
		for i, w := range worksSeed {
			secID, ok := sectionIDs[w.section]
			if !ok {
				continue
			}
			r := core.NewRecord(worksCol)
			r.Set("section", secID)
			r.Set("title", w.title)
			r.Set("slug", w.slug)
			r.Set("year", w.year)
			r.Set("role", w.role)
			r.Set("description", w.description)
			r.Set("credits", w.credits)
			r.Set("sort_order", i+1)
			r.Set("active", true)
			r.Set("featured", w.featured)
			if err := app.Save(r); err != nil {
				return err
			}
			for j, l := range w.links {
				lr := core.NewRecord(workLinksCol)
				lr.Set("work", r.Id)
				lr.Set("label", l.label)
				lr.Set("url", l.url)
				lr.Set("kind", l.kind)
				lr.Set("sort_order", j+1)
				if err := app.Save(lr); err != nil {
					return err
				}
			}
		}

		// ---------- press ----------
		pressCol, err := app.FindCollectionByNameOrId("press")
		if err != nil {
			return err
		}
		for i, p := range pressSeed {
			r := core.NewRecord(pressCol)
			r.Set("title", p.title)
			r.Set("publication", p.publication)
			r.Set("url", p.url)
			r.Set("date", p.date)
			r.Set("excerpt", p.excerpt)
			r.Set("sort_order", i+1)
			r.Set("active", true)
			if err := app.Save(r); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		// Down: borrar registros (no las colecciones).
		for _, name := range []string{"press", "work_links", "works", "sections", "site_settings"} {
			col, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				continue
			}
			records, err := app.FindAllRecords(col.Name)
			if err != nil {
				continue
			}
			for _, r := range records {
				_ = app.Delete(r)
			}
		}
		return nil
	})
}

type sectionSeed struct {
	name, slug, desc string
}

var sectionsSeed = []sectionSeed{
	{"Performance", "performance", ""},
	{"Teatro", "teatro", ""},
	{"Cine", "cine", ""},
	{"Videoclips", "videoclips", ""},
	{"Talleres", "talleres", ""},
	{"Salud Mental", "salud-mental", ""},
}

type linkSeed struct {
	label, url, kind string
}

type workSeed struct {
	section     string
	title       string
	slug        string
	year        string
	role        string
	description string
	credits     string
	featured    bool
	links       []linkSeed
}

var worksSeed = []workSeed{
	// ---------- PERFORMANCE ----------
	{
		section:  "performance",
		title:    "BZD el entierro",
		slug:     "bzd-el-entierro",
		year:     "2024",
		role:     "Creadora e intérprete",
		featured: true,
		description: "Pieza performática estrenada en 2024 en el Museo Nacional Reina Sofía, " +
			"dentro de la residencia artística de Tejidos Conjuntivos bajo la dirección de Germán Labrador. " +
			"Combina instalación escénica, performance site-specific y cine de 16 mm.",
		links: []linkSeed{
			{"Ver trailer", "https://vimeo.com/bzd-el-entierro-trailer", "trailer"},
			{"Ver obra completa", "https://vimeo.com/bzd-el-entierro", "full_work"},
		},
	},
	{
		section: "performance",
		title:   "BZD (benzodiacepinas) HASTA LOS HUESOS",
		slug:    "bzd-hasta-los-huesos",
		year:    "2023",
		role:    "Creadora e intérprete",
		description: "Estrenada en diciembre 2023 en el Festival Subterráneo Escénico " +
			"(IBERESCENA) de Ciudad de México.",
		links: []linkSeed{
			{"Ver dossier", "https://drive.google.com/bzd-dossier", "dossier"},
		},
	},
	{
		section:     "performance",
		title:       "PICNIC 2023",
		slug:        "picnic-2023",
		year:        "2023",
		role:        "Performer",
		description: "Performance en el evento PICNIC del Museo Nacional Reina Sofía.",
		links: []linkSeed{
			{"Ver video", "https://vimeo.com/picnic-2023", "full_work"},
			{"Ver certificado", "https://museoreinasofia.es/picnic-2023", "external"},
		},
	},
	{
		section: "performance",
		title:   "Artificio para atravesar la psicosis",
		slug:    "artificio-psicosis",
		year:    "2023",
		role:    "Dramaturgia y dirección",
		description: "Pieza escénica estrenada en mayo 2023 en el Auditorio Sabatini " +
			"del Museo Nacional Reina Sofía. Presentada también en el festival de teatro de Tara " +
			"en las Islas Canarias (mayo 2023) con una conferencia performática a partir de documentos " +
			"de la cárcel de Ezeiza.",
		links: []linkSeed{
			{"Ver dossier", "https://drive.google.com/artificio-dossier", "dossier"},
		},
	},
	{
		section:     "performance",
		title:       "The Castle of Crossed Destinies",
		slug:        "castle-crossed-destinies",
		year:        "2022",
		role:        "Performer",
		description: "Performer de la instalación de Leonor Serrano Rivas en el Museo Nacional Reina Sofía (septiembre 2022).",
		links: []linkSeed{
			{"Link a evento", "https://museoreinasofia.es/castle-crossed-destinies", "external"},
		},
	},
	{
		section:     "performance",
		title:       "Artificio para atravesar la psicosis y los cuerpos dóciles",
		slug:        "artificio-cuerpos-dociles",
		year:        "2022",
		role:        "Dramaturgia y dirección",
		description: "Apertura del proceso de creación del Máster en Práctica Escénica y Cultura Visual del Museo Reina Sofía y la Universidad de Castilla-La Mancha (octubre 2022).",
		links: []linkSeed{
			{"Link de trailer", "https://vimeo.com/artificio-trailer", "trailer"},
			{"Link de material documental", "https://drive.google.com/artificio-docs", "dossier"},
		},
	},
	{
		section:     "performance",
		title:       "Stop exclusión sanitaria",
		slug:        "stop-exclusion-sanitaria",
		year:        "2022",
		role:        "Creadora y directora",
		description: "Performance realizada en el evento PICNIC del Museo Nacional Reina Sofía (junio 2022) por Museo Situado.",
		links: []linkSeed{
			{"Link a evento", "https://museoreinasofia.es/stop-exclusion-sanitaria", "external"},
		},
	},

	// ---------- TEATRO ----------
	{
		section:  "teatro",
		title:    "Vitis Vinífera",
		slug:     "vitis-vinifera",
		year:     "2025",
		role:     "Creadora y directora",
		featured: true,
		description: "“Dicen que siempre hay algo de locura en el amor. " +
			"Pero siempre hay algo de razón en la locura. " +
			"Dicen que siempre hay algo de locura en el amor, pero están equivocados. " +
			"No hay algo de locura, no es una parte, es todo. El amor es locura o no es nada.”\n\n" +
			"Estrenada como creadora en el marco del festival SURGE de Madrid en otoño 2025; " +
			"participa del ciclo de Teatro Argentino del Umbral de Primavera en febrero 2026.",
		credits: "Actúa: Oscar Bell · Dirección: Amelia Repetto · Coreografía: Nina Gorostiza · Luces y Mapping: Werner Faramarz · Diseño de iluminación: Belén Abarza",
		links: []linkSeed{
			{"Trailer", "https://vimeo.com/vitis-vinifera-trailer", "trailer"},
			{"Obra Completa", "https://vimeo.com/vitis-vinifera", "full_work"},
			{"Dossier", "https://drive.google.com/vitis-dossier", "dossier"},
		},
	},
	{
		section: "teatro",
		title:   "Romance de la Negra Rubia",
		slug:    "romance-negra-rubia",
		year:    "2024",
		role:    "Dirección",
		description: "Estrenada el 4 de febrero de 2024 en el Ciclo de Teatro Argentino del Umbral de Primavera Madrid " +
			"y presentada en mayo 2022 en el Encuentro de Teatro Argentino de Barcelona en el Teatre L'Estranger. " +
			"Residencia de creación NUTRIENTES en teatro Pradillo (septiembre 2023) con la obra de Gabriela Cabezón Cámara, " +
			"dirección de Emilia Dulom. Gira por La Rioja (España) en abril 2025: Logroño, Ábalos, Leza, Peroblasco.",
		links: []linkSeed{
			{"Ver dossier", "https://drive.google.com/romance-dossier", "dossier"},
			{"Ver trailer", "https://vimeo.com/romance-negra-rubia", "trailer"},
		},
	},
	{
		section:     "teatro",
		title:       "Instalación performática de teatro y fotografía en espacio Ballux",
		slug:        "ballux-instalacion",
		year:        "2021",
		role:        "Performer",
		description: "Instalación performática en el espacio Ballux.",
		credits:     "PH: Anabella Sarrias · Espacio Ballux",
	},
	{
		section: "teatro",
		title:   "La débil mental",
		slug:    "la-debil-mental",
		year:    "2021",
		role:    "Asistencia artística",
		description: "Estreno 14 de diciembre 2021. Obra en proceso ganadora de Iberescena, sobre la novela " +
			"“La débil mental” de Ariana Harwicz.",
		credits: "Teatro: Galpón de Guevara · Dramaturgia: Ariana Harwicz · Dirección: Cristina Banegas",
		links: []linkSeed{
			{"Suplemento Cultural Ñ Clarín", "https://www.clarin.com/cultura/la-debil-mental.html", "press"},
			{"Trabajos previos sobre la novela", "https://www.ejemplo.com/debil-mental-trabajos", "external"},
			{"Instagram @ladebilmentalteatro", "https://instagram.com/ladebilmentalteatro", "instagram"},
		},
	},
	{
		section: "teatro",
		title:   "Las cuerdas",
		slug:    "las-cuerdas",
		year:    "2019-2021",
		role:    "Actriz",
		description: "Obra ganadora del concurso Óperas Primas del Centro Cultural Ricardo Rojas. " +
			"Estrenada en 2019 en el Centro Cultural Rojas; funciones en noviembre y diciembre 2020 en el Teatro Espacio Callejón; " +
			"temporada 2021 en el Espacio Callejón.",
		credits: "Dramaturgia y dirección: Ana Schimelman · Co-protagonista: Fiamma Carranza Macchi",
		links: []linkSeed{
			{"Ficha de la obra", "https://alternativateatral.com/las-cuerdas", "external"},
			{"Nota Página 12 Suplemento Radar", "https://pagina12.com.ar/radar-las-cuerdas", "press"},
			{"Nota La Nación", "https://lanacion.com.ar/las-cuerdas", "press"},
			{"Nota Página 12", "https://pagina12.com.ar/las-cuerdas", "press"},
			{"Nota Fervor", "https://fervor.com.ar/las-cuerdas", "press"},
		},
	},
	{
		section:     "teatro",
		title:       "El péndulo",
		slug:        "el-pendulo",
		year:        "2019",
		role:        "Actriz protagónica",
		description: "Obra seleccionada para participar del Festival Fauna.",
		credits:     "Dramaturgia y dirección: Magdalena Yomha · Teatro: Centro Cultural General San Martín (CCGSM)",
		links: []linkSeed{
			{"Festival Fauna", "https://festivalfauna.com.ar/el-pendulo", "external"},
		},
	},
	{
		section: "teatro",
		title:   "Lo que se dice",
		slug:    "lo-que-se-dice",
		year:    "2018",
		role:    "Actriz",
	},
	{
		section:     "teatro",
		title:       "Amor, mentiras y dinero",
		slug:        "amor-mentiras-dinero",
		year:        "2018",
		role:        "Actriz protagónica",
		description: "Microteatro BA.",
		credits:     "Dramaturgia y dirección: Eugenio Soto",
		links: []linkSeed{
			{"Ficha (alternativa teatral)", "https://alternativateatral.com/amor-mentiras-dinero", "external"},
		},
	},
	{
		section:     "teatro",
		title:       "Ensueños",
		slug:        "ensuenos",
		year:        "2018",
		role:        "Actriz protagónica",
		description: "Microteatro BA.",
		credits:     "Dramaturgia: Laura Manson y Lucila Brea · Dirección: Laura Manson",
		links: []linkSeed{
			{"Ficha (alternativa teatral)", "https://alternativateatral.com/ensuenos", "external"},
		},
	},
	{
		section: "teatro",
		title:   "El tilo",
		slug:    "el-tilo",
		year:    "2017",
		role:    "Actriz",
		credits: "Dramaturgia: Agustín Maradei",
	},
	{
		section:     "teatro",
		title:       "La liebre y la tortuga",
		slug:        "la-liebre-y-la-tortuga",
		year:        "2017",
		role:        "Diseño de vestuario",
		description: "Teatro Nacional Cervantes, Laboratorio de Creación I.",
		credits:     "Dramaturgia y dirección: Ricardo Bartís",
		links: []linkSeed{
			{"Ficha de la obra (Cervantes)", "https://teatrocervantes.gob.ar/la-liebre-y-la-tortuga", "external"},
			{"Nota diario Página 12", "https://pagina12.com.ar/la-liebre-y-la-tortuga", "press"},
			{"Nota diario Clarín", "https://clarin.com/la-liebre-y-la-tortuga", "press"},
		},
	},
	{
		section:     "teatro",
		title:       "La Pirámide",
		slug:        "la-piramide",
		year:        "2016",
		role:        "Actriz protagónica",
		description: "Teatro Dínamo.",
		credits:     "Dramaturgia: Copi · Dirección: Daniela Regert",
	},
	{
		section: "teatro",
		title:   "Zoom in 90s",
		slug:    "zoom-in-90s",
		year:    "2016",
		role:    "Actriz",
		credits: "Dirección: Rubén Sabadini",
	},
	{
		section:     "teatro",
		title:       "Las multitudes",
		slug:        "las-multitudes",
		year:        "2013",
		role:        "Actriz",
		description: "Centro Cultural General San Martín.",
		credits:     "Dramaturgia y dirección: Federico León",
		links: []linkSeed{
			{"Ficha (alternativa teatral)", "https://alternativateatral.com/las-multitudes", "external"},
			{"Trailer", "https://vimeo.com/las-multitudes", "trailer"},
			{"Nota diario Página 12", "https://pagina12.com.ar/las-multitudes", "press"},
			{"Nota diario La Nación", "https://lanacion.com.ar/las-multitudes", "press"},
			{"Otras notas 1", "https://otro.com/las-multitudes-1", "press"},
			{"Otras notas 2", "https://otro.com/las-multitudes-2", "press"},
		},
	},
	{
		section: "teatro",
		title:   "Cinefilia",
		slug:    "cinefilia",
		year:    "2013-2012",
		role:    "Actriz",
		credits: "Dirección: Aníbal Gulluni",
	},
	{
		section: "teatro",
		title:   "Víspera de elecciones",
		slug:    "vispera-de-elecciones",
		year:    "2011",
		role:    "Actriz",
		credits: "Óperas Primas del Centro Cultural Ricardo Rojas",
	},
	{
		section: "teatro",
		title:   "El juicio de Lady Macbeth",
		slug:    "juicio-lady-macbeth",
		year:    "2010-2009",
		role:    "Actriz",
	},

	// ---------- CINE ----------
	{
		section:     "cine",
		title:       "Inconsciente colectivo",
		slug:        "inconsciente-colectivo",
		year:        "",
		role:        "Actriz",
		description: "Serie Amazon Prime Video.",
		credits:     "Dirección: Mariano Hueter",
		links: []linkSeed{
			{"Ficha Prime Video", "https://primevideo.com/inconsciente-colectivo", "external"},
		},
	},
	{
		section:     "cine",
		title:       "Noemí Gold",
		slug:        "noemi-gold",
		year:        "",
		role:        "Actriz",
		description: "Largometraje Amazon Prime Video.",
		credits:     "Dirección: Dan Rubenstein",
		links: []linkSeed{
			{"Trailer", "https://vimeo.com/noemi-gold-trailer", "trailer"},
			{"Cineramaplus", "https://cineramaplus.com/noemi-gold", "press"},
			{"Rogers Movie Nation", "https://rogersmovienation.com/noemi-gold", "press"},
			{"Vice.com", "https://vice.com/noemi-gold", "press"},
			{"Cinéfiloserial.com.ar", "https://cinefiloserial.com.ar/noemi-gold", "press"},
			{"Mubi.com", "https://mubi.com/films/noemi-gold", "external"},
			{"Imdb.com", "https://imdb.com/title/noemi-gold", "external"},
			{"Filmaffinity.com", "https://filmaffinity.com/noemi-gold", "external"},
			{"Topic.com", "https://topic.com/noemi-gold", "external"},
			{"Deadline.com", "https://deadline.com/noemi-gold", "press"},
		},
	},
	{
		section:     "cine",
		title:       "Crónicas ferreteras",
		slug:        "cronicas-ferreteras",
		year:        "",
		role:        "Actriz",
		description: "Serie en Cinearplay. Episodio: Síndrome de Estocolmo.",
		credits:     "Dirección: Mariano Fernández",
		links: []linkSeed{
			{"Cinear", "https://cine.ar/cronicas-ferreteras", "external"},
			{"Episodio en YouTube", "https://youtube.com/cronicas-ferreteras-estocolmo", "full_work"},
		},
	},
	{
		section:     "cine",
		title:       "Todo lo que veo es mío",
		slug:        "todo-lo-que-veo-es-mio",
		year:        "",
		role:        "Actriz bolo",
		description: "Largometraje.",
		credits:     "Dirección: Mariano Galperín y Román Podolsky",
		links: []linkSeed{
			{"Trailer", "https://vimeo.com/todo-lo-que-veo-es-mio", "trailer"},
		},
	},
	{
		section:     "cine",
		title:       "El éxtasis de Santa Teresa",
		slug:        "extasis-santa-teresa",
		year:        "2021",
		role:        "Actriz",
		description: "Mediometraje.",
		credits:     "Dirección: Lucas Matranga",
	},
	{
		section:     "cine",
		title:       "Veredas",
		slug:        "veredas",
		year:        "",
		role:        "Actriz bolo",
		description: "Cinearplay.",
		credits:     "Dirección: Fernando Cricenti",
	},

	// ---------- VIDEOCLIPS ----------
	{
		section: "videoclips",
		title:   "No todo es color de rosa",
		slug:    "no-todo-es-color-de-rosa",
		year:    "",
		role:    "Actriz",
		credits: "de Nahuel Briones",
	},
	{
		section: "videoclips",
		title:   "Casa roja",
		slug:    "casa-roja",
		year:    "",
		role:    "Actriz",
		credits: "de Pil y los Violadores",
	},
	{
		section: "videoclips",
		title:   "Sensaciones",
		slug:    "sensaciones",
		year:    "",
		role:    "Actriz",
		credits: "de Arbolito",
	},
	{
		section: "videoclips",
		title:   "Linda",
		slug:    "linda",
		year:    "",
		role:    "Actriz",
		credits: "de Marcelo Ezquiaga",
	},

	// ---------- TALLERES ----------
	{
		section:     "talleres",
		title:       "Taller en La Quimera",
		slug:        "taller-la-quimera",
		year:        "2022",
		role:        "Docente",
		description: "Taller de teatro mediadora cultural con población afroespañola en el barrio de Lavapiés. Invitada por la asociación Beshawear en el espacio okupa “La Quimera” (agosto 2022).",
	},

	// ---------- SALUD MENTAL ----------
	{
		section:  "salud-mental",
		title:    "Intervenciones performáticas de arte y salud mental en el penal de Máxima Seguridad de Ezeiza",
		slug:     "ezeiza-intervenciones",
		year:     "",
		role:     "Creadora y directora",
		featured: true,
		description: "Videos de trabajos dentro de la cárcel, pabellón para pacientes con trastornos de salud mental privados de su libertad.",
		links: []linkSeed{
			{"Videoclip “Trato por igual”", "https://youtube.com/trato-por-igual", "full_work"},
			{"Clase de teatro en penal de Ezeiza", "https://youtube.com/clase-teatro-ezeiza", "full_work"},
		},
	},
	{
		section:     "salud-mental",
		title:       "El juicio F22.0",
		slug:        "el-juicio-f22",
		year:        "2010",
		role:        "Productora y directora de actores",
		description: "Cortometraje sobre trastorno Delirante crónico. Realizado en el marco de la residencia de Salud Mental del Hospital Pirovano (Buenos Aires) y el Centro de Salud Mental N1 Dr. Hugo Rosarios. Proyectado en el Congreso internacional de Psiquiatría APSA (Mar del Plata, 2010) y en las jornadas de residentes de salud mental.",
		links: []linkSeed{
			{"Link YouTube", "https://youtube.com/el-juicio-f22", "full_work"},
		},
	},
}

type pressSeedItem struct {
	title, publication, url, date, excerpt string
}

var pressSeed = []pressSeedItem{
	{
		title:       "De Buenos Aires a Madrid sin escalas: la psiquiatra argentina que abandonó una profesión de 11 años por amor al teatro",
		publication: "El mundo",
		url:         "https://elmundo.es/amelia-repetto-psiquiatra-teatro",
		date:        "2023-09-04",
		excerpt:     "La psiquiatra que dejó todo por el teatro: «Dio sentido a mi existencia».",
	},
	{
		title:       "Amelia Repetto presentó “Artificio para atravesar la psicosis y los cuerpos dóciles” en la Asociación ATLAS",
		publication: "Masescena",
		url:         "https://masescena.com/amelia-repetto-artificio-atlas",
		date:        "2023-05-23",
		excerpt:     "Presentación con participación ciudadana dentro del Festival Tara.",
	},
}

const bioES = `<p>Actriz, directora de teatro y psiquiatra especializada en arte y salud mental. Miembro de ARTEA (<a href="https://artea.uclm.es/">artea.uclm.es</a>).</p>

<p>Se formó como actriz en la Universidad Nacional de las Artes (UNA), departamento de Dramáticas.</p>

<p>En septiembre de 2025 estrenó como creadora <em>VITIS VINIFERA</em> en el marco del festival SURGE de Madrid en Otoño; dicha obra participó del ciclo de Teatro Argentino del Umbral de primavera, febrero 2026. Participó como actriz en Micro Teatro por Dinero en la obra <em>Paraíso apartment</em> junto a Antonella Mazota y la dirección de Leo Bartolotta. La obra fue seleccionada por Microteatro para participar de los siguientes festivales: Benidorm, Lazarillo, Centro de Creación Cárcel de Segovia. Estará entre octubre y noviembre del presente año en la Micro selección en la sección clásicos.</p>

<p>En marzo 2025 participó de la obra <em>LAVAPIÉS</em> dirigida por Fernando Ferrer en el Teatro del Barrio (premio nacional del teatro 2024).</p>

<p>En abril 2025 realizó gira por la provincia de La Rioja España con el <em>Romance de la Negra Rubia</em>, por las ciudades de Logroño, Ábalos, Leza y Peroblasco.</p>

<p>Trabaja realizando laboratorios de artes escénicas en el Centro de Arte 2 de Mayo. De octubre de 2024 a septiembre 2025 trabaja como performer de la obra de Santiago Sierra y dicta laboratorios de artes escénicas en el Museo CA2M (Centro de Arte Dos de Mayo).</p>

<p>2024 artista residente de Tejidos Conjuntivos del Museo Reina Sofía, bajo la dirección de German Labrador. Espacio en el cual investiga una pieza que combina una instalación escénica, una performance site specific y cine de 16 mm. Marco en el cual estrenó BZD EL ENTIERRO.</p>

<p>Realizó en 2022 el Máster en Práctica Escénica y Cultura Visual del Museo Nacional Reina Sofía y la Universidad de Castilla la Mancha.</p>

<p>Complementó su formación actoral con Carlota Ferrer (Madrid), Claudia Cantero, Ricardo Bartis, Federico León, Alejandro Catalán, Guillermo Angeleli, Alfredo Castro (Chile) y Lito Cruz.</p>

<p>Fue una de las cuatro artistas seleccionadas para participar de la residencia de lectura de Dorothy Michaels en octubre de 2023 bajo la coordinación de María Jerez.</p>

<p>Participó de la edición XI del FESTIVAL SUBTERRÁNEO DE MÉXICO (2023) junto a IBERESCENA, a desarrollarse en DF MÉXICO del 1 al 9 de diciembre de 2023 como creadora, en donde estrenó una pieza inédita creada por el festival <em>BZD hasta los huesos</em>.</p>

<p>Participó de la residencia de septiembre de 2023 de la residencia NUTRIENTES en teatro Pradillo con la obra <em>ROMANCE DE LA NEGRA RUBIA</em>, obra argentina de la autora Gabriela Cabezón Cámara, bajo la dirección de Emilia Dulom. Esta obra se estrenó el 4 de febrero de 2024 en el Ciclo de Teatro Argentino del Umbral de Primavera en Madrid y participó del Encuentro de Teatro Argentino de Barcelona en el Teatre L'Estranger en mayo de 2022.</p>

<p>Realizó como actriz, directora y dramaturga un proyecto escénico, llamado <em>Artificio para atravesar la psicosis</em>, estrenado en el auditorio Sabatini del Museo Reina Sofía y desarrollado a lo largo de diferentes locaciones del museo. Participó con esta pieza en el festival de teatro de Tara en las Islas Canarias en el mes de mayo de 2023. Esta pieza consistía en una conferencia performática realizada con documentos tomados en la cárcel de Ezeiza y una performance.</p>

<p>En febrero de 2023 estrenó la obra <em>Las Cuerdas</em> de Ana Schimelman en el teatro Umbral de Primavera en el Festival de Teatro Argentino, obra portada de la revista de teatro Godot de febrero.</p>

<p>En junio de 2023 trabajó en la dirección creativa del evento <em>Picnic</em>, tanto en la creación de los módulos escenográficos como en la coordinación de la performance para la visibilización de las campañas artístico-políticas del evento, en sintonía con las líneas curatoriales de la exposición <em>Maquinaciones</em> (2023) del Museo Reina Sofía, comisariado por Pablo Allepuz, Manuel Borja-Villel, Ileana Fokianaki, Rafael García y Teresa Velázquez. Participó como performer en septiembre de 2022 en la instalación <em>The Castle of Crossed Destinies</em> de la artista contemporánea Leonor Serrano Rivas en el Museo Nacional Reina Sofía. También colaboró con Museo Situado en el diseño y la dirección de la performance <em>Stop Exclusión Sanitaria</em> en el marco del evento Picnic del Museo Nacional Reina Sofía.</p>

<p>Trabajó en diversos proyectos teatrales como <em>Las Cuerdas</em> (Ana Schimelman), <em>El Péndulo</em> (Magdalena Yomha), <em>Las Multitudes</em> (Federico León), <em>Víspera de elecciones</em> del concurso Óperas Primas del Centro Cultural Ricardo Rojas, estrenada en febrero de 2021 en el Espacio Callejón; <em>El Péndulo</em> en el Camarín de las Musas, <em>Cinefilia</em> de Aníbal Gulluni, en la Carpintería Teatro. <em>Amor, Mentiras y Dinero</em> de Eugenio Soto, en Microteatro; <em>Ensueños</em> de Laura Manson, en Microteatro; <em>Lo que se dice</em> de Andrés Raiano, en Pelonia; <em>El Tilo</em> de Agustín Maradei, en Laboratorio Teatral; <em>Desde Vera 108</em> de Rubén Sabadini, y <em>Zoom in 90s</em> de Candelaria Sabagh.</p>

<p>Como docente se desempeñó dictando talleres de actuación y laboratorios de creación en el Umbral de Primavera Madrid, La Parcería, y Centro de Arte Dos de Mayo (CA2M). En Argentina (2020-2021), se desempeñó como teaching assistant en la Universidad Nacional de las Artes (UNA), Departamento de Artes Visuales, en Acting Direction I y II (Costa Costa). También se desempeñó como asistente I en Acting (Bernardo Cappa) en el Departamento de Artes Dramáticas de UNA (2011-2013).</p>`

const bioEN = `<p>Actress, theatre director, and psychiatrist specialized in art and mental health. Member of ARTEA (<a href="https://artea.uclm.es/">artea.uclm.es</a>).</p>

<p>She trained as an actress at the National University of the Arts (UNA), Department of Dramatic Arts.</p>

<p>In September 2025, she premiered her original creation <em>VITIS VINIFERA</em> as part of the SURGE Festival in Madrid (Autumn edition). The piece later participated in the Argentine Theatre Cycle at Umbral de Primavera in February 2026. She performed in Microteatro por Dinero in the play <em>Paraíso Apartment</em> alongside Antonella Mazota, directed by Leo Bartolotta. The production was selected by Microteatro to participate in several festivals, including Benidorm, Lazarillo, and the Centro de Creación Cárcel de Segovia.</p>

<p>In March 2025, she performed in <em>LAVAPIÉS</em>, directed by Fernando Ferrer at Teatro del Barrio (National Theatre Award 2024).</p>

<p>In April 2025, she toured La Rioja (Spain) with <em>Romance de la Negra Rubia</em>, performing in Logroño, Ábalos, Leza, and Peroblasco.</p>

<p>She currently leads performing arts laboratories at Centro de Arte Dos de Mayo (CA2M). From October 2024 to September 2025, she worked as a performer in a piece by Santiago Sierra and also taught performing arts laboratories at CA2M.</p>

<p>In 2024, she was an artist-in-residence at Tejidos Conjuntivos at Museo Reina Sofía, under the direction of Germán Labrador. There, she developed a project combining scenic installation, site-specific performance, and 16mm film, within which she premiered <em>BZD EL ENTIERRO</em>.</p>

<p>In 2022, she completed the Master's Degree in Performing Practice and Visual Culture at Museo Reina Sofía and the University of Castilla-La Mancha.</p>

<p>She complemented her acting training with Carlota Ferrer (Madrid), Claudia Cantero, Ricardo Bartis, Federico León, Alejandro Catalán, Guillermo Angeleli, Alfredo Castro (Chile), and Lito Cruz.</p>

<p>She was one of four artists selected for the Dorothy Michaels Reading Residency (October 2023), coordinated by María Jerez. She participated in the XI Subterraneo Festival in Mexico (2023), held in Mexico City from December 1-9, where she premiered the original work <em>BZD hasta los huesos</em>.</p>

<p>In September 2023, she took part in the NUTRIENTES residency at Teatro Pradillo with <em>ROMANCE OF THE NEGRA RUBIA</em> by Argentine author Gabriela Cabezón Cámara. The piece premiered on February 4, 2024, at the Argentine Theatre Cycle (Umbral de Primavera), and was later presented at the Argentine Theatre Meeting in Barcelona at Teatre L'Estranger (May 2022).</p>

<p>She created, directed, and performed <em>Artificio para atravesar la psicosis</em>, a scenic project premiered at the Sabatini Auditorium of Museo Reina Sofía and developed across various locations within the museum. The piece was also presented at the Tara Theatre Festival (Canary Islands, May 2023). In February 2023, she premiered <em>Las Cuerdas</em> by Ana Schimelman at Umbral de Primavera during the Argentine Theatre Festival. The production was featured on the cover of Revista Godot (February issue).</p>

<p>In June 2023, she worked on the creative direction of the event <em>Picnic</em>, contributing to scenographic modules and coordinating performance actions aligned with the artistic-political campaigns of the exhibition <em>Maquinaciones</em> (Museo Reina Sofía), curated by Pablo Allepuz, Manuel Borja-Villel, Ileana Fokianaki, Rafael García, and Teresa Velázquez. In September 2022, she participated as a performer in <em>THE CASTLE OF CROSSED DESTINIES</em> by contemporary artist Leonor Serrano Rivas at Museo Reina Sofía. She also collaborated with Museo Situado in designing and directing the performance <em>Stop Exclusión Sanitaria</em> at the Picnic event.</p>

<p>She has performed in numerous theatre productions, including: <em>Las Cuerdas</em> (Ana Schimelman), <em>El Péndulo</em> (Magdalena Yomha), <em>Las Multitudes</em> (Federico León), <em>Cinefilia</em> (Aníbal Gulluni), <em>Amor, Mentiras y Dinero</em> (Eugenio Soto), <em>Ensueños</em> (Laura Manson), <em>Lo que se dice</em> (Andrés Raiano), <em>El Tilo</em> (Agustín Maradei), <em>Desde Vera 108</em> (Rubén Sabadini), and <em>Zoom in 90s</em> (Candelaria Sabagh).</p>

<p>As a teacher, she currently leads acting workshops and creative laboratories at Umbral de Primavera (Madrid), La Parcería, and Centro de Arte Dos de Mayo (CA2M). In Argentina (2020-2021), she worked as a teaching assistant at the National University of the Arts (UNA), Department of Visual Arts, in Acting Direction I & II (Costa Costa). She also served as assistant I in Acting (Bernardo Cappa) at UNA's Department of Dramatic Arts (2011-2013).</p>`
