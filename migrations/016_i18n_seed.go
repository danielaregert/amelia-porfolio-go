package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Carga las traducciones al inglés de todo el contenido seedeado.
// bio_en ya fue cargada en la migración 002.
func init() {
	m.Register(func(app core.App) error {

		// ---------- site_settings ----------
		if rec, err := app.FindFirstRecordByFilter("site_settings", "id != ''"); err == nil && rec != nil {
			rec.Set("tagline_en", "artistic portfolio")
			_ = app.Save(rec)
		}

		// ---------- sections ----------
		sectionEN := map[string]string{
			"performance":  "Performance",
			"teatro":       "Theatre",
			"cine":         "Film",
			"videoclips":   "Music videos",
			"talleres":     "Workshops",
			"salud-mental": "Mental Health",
		}
		for slug, en := range sectionEN {
			if s, err := app.FindFirstRecordByFilter("sections", "slug = {:s}",
				map[string]any{"s": slug}); err == nil && s != nil {
				s.Set("name_en", en)
				_ = app.Save(s)
			}
		}

		// ---------- works ----------
		for slug, tr := range worksEN {
			if w, err := app.FindFirstRecordByFilter("works", "slug = {:s}",
				map[string]any{"s": slug}); err == nil && w != nil {
				if tr.title != "" {
					w.Set("title_en", tr.title)
				}
				if tr.role != "" {
					w.Set("role_en", tr.role)
				}
				if tr.description != "" {
					w.Set("description_en", tr.description)
				}
				if tr.credits != "" {
					w.Set("credits_en", tr.credits)
				}
				_ = app.Save(w)
			}
		}

		// ---------- press ----------
		for _, p := range pressEN {
			if rec, err := app.FindFirstRecordByFilter("press", "title = {:t}",
				map[string]any{"t": p.titleES}); err == nil && rec != nil {
				rec.Set("title_en", p.titleEN)
				rec.Set("excerpt_en", p.excerptEN)
				_ = app.Save(rec)
			}
		}

		return nil
	}, func(app core.App) error { return nil })
}

type workTranslation struct {
	title, role, description, credits string
}

var worksEN = map[string]workTranslation{
	// ---------- PERFORMANCE ----------
	"bzd-el-entierro": {
		title: "BZD the burial",
		role:  "Creator and performer",
		description: "Performance piece premiered in 2024 at the Museo Nacional Reina Sofía, " +
			"as part of the Tejidos Conjuntivos artist residency under the direction of Germán Labrador. " +
			"Combines scenic installation, site-specific performance and 16mm film.",
	},
	"bzd-hasta-los-huesos": {
		title:       "BZD (benzodiazepines) TO THE BONE",
		role:        "Creator and performer",
		description: "Premiered in December 2023 at the Subterráneo Escénico Festival (IBERESCENA) in Mexico City.",
	},
	"picnic-2023": {
		title:       "PICNIC 2023",
		role:        "Performer",
		description: "Performance at the PICNIC event of Museo Nacional Reina Sofía.",
	},
	"artificio-psicosis": {
		title: "Device for traversing psychosis",
		role:  "Playwright and director",
		description: "Scenic piece premiered in May 2023 at the Sabatini Auditorium of " +
			"Museo Nacional Reina Sofía. Also presented at the Tara Theatre Festival in the Canary Islands " +
			"(May 2023) with a performative lecture based on documents from Ezeiza prison.",
	},
	"castle-crossed-destinies": {
		title:       "The Castle of Crossed Destinies",
		role:        "Performer",
		description: "Performer in the installation by Leonor Serrano Rivas at Museo Nacional Reina Sofía (September 2022).",
	},
	"artificio-cuerpos-dociles": {
		title:       "Device for traversing psychosis and docile bodies",
		role:        "Playwright and director",
		description: "Open creation process of the Master's Program in Scenic Practice and Visual Culture at Museo Reina Sofía and the University of Castilla-La Mancha (October 2022).",
	},
	"stop-exclusion-sanitaria": {
		title:       "Stop health exclusion",
		role:        "Creator and director",
		description: "Performance held at the PICNIC event of Museo Nacional Reina Sofía (June 2022), organized by Museo Situado.",
	},

	// ---------- TEATRO ----------
	"vitis-vinifera": {
		title: "Vitis Vinifera",
		role:  "Creator and director",
		description: `"They say there is always some madness in love. ` +
			`But there is also some reason in madness. ` +
			`They say there is always some madness in love, but they are wrong. ` +
			`There is no 'some' madness, it is not a part, it is everything. Love is madness or nothing at all."` +
			"\n\nPremiered as creator at the SURGE Festival in Madrid (Autumn 2025); " +
			"part of the Argentine Theatre Cycle at Umbral de Primavera in February 2026.",
		credits: "Acting: Oscar Bell · Direction: Amelia Repetto · Choreography: Nina Gorostiza · Lights & Mapping: Werner Faramarz · Lighting Design: Belén Abarza",
	},
	"romance-negra-rubia": {
		title: "Romance of the Blonde Black Woman",
		role:  "Direction",
		description: "Premiered on February 4, 2024 at the Argentine Theatre Cycle of Umbral de Primavera Madrid, " +
			"and presented in May 2022 at the Argentine Theatre Meeting in Barcelona at Teatre L'Estranger. " +
			"NUTRIENTES residency at Teatro Pradillo (September 2023) with Gabriela Cabezón Cámara's work, directed by Emilia Dulom. " +
			"Tour through La Rioja (Spain) in April 2025: Logroño, Ábalos, Leza, Peroblasco.",
	},
	"ballux-instalacion": {
		title:       "Performative installation of theatre and photography at Ballux space",
		role:        "Performer",
		description: "Performative installation at Ballux space.",
		credits:     "PH: Anabella Sarrias · Ballux Space",
	},
	"la-debil-mental": {
		title:       "The mentally weak",
		role:        "Artistic assistance",
		description: `Premiered December 14, 2021. Work-in-progress, Iberescena award winner, based on Ariana Harwicz's novel "La débil mental".`,
		credits:     "Venue: Galpón de Guevara · Playwriting: Ariana Harwicz · Direction: Cristina Banegas",
	},
	"las-cuerdas": {
		title: "The strings",
		role:  "Actress",
		description: "Winner of the First Works Competition at Centro Cultural Ricardo Rojas. " +
			"Premiered in 2019 at Centro Cultural Rojas; shows in November and December 2020 at Teatro Espacio Callejón; " +
			"2021 season at Espacio Callejón.",
		credits: "Playwriting and direction: Ana Schimelman · Co-lead: Fiamma Carranza Macchi",
	},
	"el-pendulo": {
		title:       "The pendulum",
		role:        "Lead actress",
		description: "Selected to participate in the Fauna Festival.",
		credits:     "Playwriting and direction: Magdalena Yomha · Venue: Centro Cultural General San Martín (CCGSM)",
	},
	"lo-que-se-dice": {
		title: "What is said",
		role:  "Actress",
	},
	"amor-mentiras-dinero": {
		title:       "Love, lies and money",
		role:        "Lead actress",
		description: "Microteatro BA.",
		credits:     "Playwriting and direction: Eugenio Soto",
	},
	"ensuenos": {
		title:       "Daydreams",
		role:        "Lead actress",
		description: "Microteatro BA.",
		credits:     "Playwriting: Laura Manson and Lucila Brea · Direction: Laura Manson",
	},
	"el-tilo": {
		title:   "The linden tree",
		role:    "Actress",
		credits: "Playwriting: Agustín Maradei",
	},
	"la-liebre-y-la-tortuga": {
		title:       "The hare and the tortoise",
		role:        "Costume design",
		description: "Teatro Nacional Cervantes, Creation Laboratory I.",
		credits:     "Playwriting and direction: Ricardo Bartís",
	},
	"la-piramide": {
		title:       "The Pyramid",
		role:        "Lead actress",
		description: "Teatro Dínamo.",
		credits:     "Playwriting: Copi · Direction: Daniela Regert",
	},
	"zoom-in-90s": {
		title:   "Zoom in 90s",
		role:    "Actress",
		credits: "Direction: Rubén Sabadini",
	},
	"las-multitudes": {
		title:       "The multitudes",
		role:        "Actress",
		description: "Centro Cultural General San Martín.",
		credits:     "Playwriting and direction: Federico León",
	},
	"cinefilia": {
		title:   "Cinephilia",
		role:    "Actress",
		credits: "Direction: Aníbal Gulluni",
	},
	"vispera-de-elecciones": {
		title:   "Election eve",
		role:    "Actress",
		credits: "First Works Competition, Centro Cultural Ricardo Rojas",
	},
	"juicio-lady-macbeth": {
		title: "The trial of Lady Macbeth",
		role:  "Actress",
	},

	// ---------- CINE ----------
	"inconsciente-colectivo": {
		title:       "Collective Unconscious",
		role:        "Actress",
		description: "Amazon Prime Video series.",
		credits:     "Direction: Mariano Hueter",
	},
	"noemi-gold": {
		title:       "Noemí Gold",
		role:        "Actress",
		description: "Amazon Prime Video feature film.",
		credits:     "Direction: Dan Rubenstein",
	},
	"cronicas-ferreteras": {
		title:       "Hardware Chronicles",
		role:        "Actress",
		description: "Series on Cinearplay. Episode: Stockholm Syndrome.",
		credits:     "Direction: Mariano Fernández",
	},
	"todo-lo-que-veo-es-mio": {
		title:       "Everything I see is mine",
		role:        "Supporting actress",
		description: "Feature film.",
		credits:     "Direction: Mariano Galperín and Román Podolsky",
	},
	"extasis-santa-teresa": {
		title:       "The ecstasy of Saint Teresa",
		role:        "Actress",
		description: "Mid-length film.",
		credits:     "Direction: Lucas Matranga",
	},
	"veredas": {
		title:       "Sidewalks",
		role:        "Supporting actress",
		description: "Cinearplay.",
		credits:     "Direction: Fernando Cricenti",
	},

	// ---------- VIDEOCLIPS ----------
	"no-todo-es-color-de-rosa": {
		title:   "Not everything is rose-colored",
		role:    "Actress",
		credits: "by Nahuel Briones",
	},
	"casa-roja": {
		title:   "Red house",
		role:    "Actress",
		credits: "by Pil y los Violadores",
	},
	"sensaciones": {
		title:   "Sensations",
		role:    "Actress",
		credits: "by Arbolito",
	},
	"linda": {
		title:   "Linda",
		role:    "Actress",
		credits: "by Marcelo Ezquiaga",
	},

	// ---------- TALLERES ----------
	"taller-la-quimera": {
		title:       "Workshop at La Quimera",
		role:        "Teacher",
		description: "Theatre workshop as cultural mediator with the Afro-Spanish community in the Lavapiés neighborhood. Invited by the Beshawear association at the squat La Quimera (August 2022).",
	},

	// ---------- SALUD MENTAL ----------
	"ezeiza-intervenciones": {
		title:       "Performative interventions of art and mental health at the Maximum Security Prison of Ezeiza",
		role:        "Creator and director",
		description: "Videos of work inside the prison, mental health ward for inmates.",
	},
	"el-juicio-f22": {
		title:       "Trial F22.0",
		role:        "Producer and acting director",
		description: "Short film on Chronic Delusional disorder. Made within the Mental Health residency at Pirovano Hospital (Buenos Aires) and the Mental Health Center N1 Dr. Hugo Rosarios. Screened at the International Psychiatry Congress APSA (Mar del Plata, 2010) and at the mental health residents' conference.",
	},
}

type pressTranslation struct {
	titleES, titleEN, excerptEN string
}

var pressEN = []pressTranslation{
	{
		titleES:   "De Buenos Aires a Madrid sin escalas: la psiquiatra argentina que abandonó una profesión de 11 años por amor al teatro",
		titleEN:   "From Buenos Aires to Madrid without stops: the Argentine psychiatrist who left an 11-year career for love of theatre",
		excerptEN: "The psychiatrist who left everything for theatre: «It gave meaning to my existence».",
	},
	{
		titleES:   "Amelia Repetto presentó “Artificio para atravesar la psicosis y los cuerpos dóciles” en la Asociación ATLAS",
		titleEN:   "Amelia Repetto presented \"Device for traversing psychosis and docile bodies\" at ATLAS Association",
		excerptEN: "Presentation with citizen participation as part of the Tara Festival.",
	},
}
