package handlers

import (
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"porfolio-amelia/internal/sanitize"
)

// safeHTML aplica sanitización en la capa de lectura como defensa en
// profundidad: aunque los datos ya se sanean al guardar, re-sanear al leer
// protege frente a edición directa de la BD o escrituras por la API de PB.
func safeHTML(s string) template.HTML {
	return template.HTML(sanitize.HTML(s))
}

var yearRe = regexp.MustCompile(`\d{4}`)

// latestYear extrae todos los años (4 dígitos) y devuelve el mayor.
// Soporta "2024", "2022-2023", "2019, 2021", "2010-2009", etc.
// Si no encuentra ninguno devuelve 0.
func latestYear(s string) int {
	matches := yearRe.FindAllString(s, -1)
	max := 0
	for _, y := range matches {
		n, err := strconv.Atoi(y)
		if err == nil && n > max {
			max = n
		}
	}
	return max
}

type pubLink struct {
	Label, URL, Kind string
}

type pubWork struct {
	Title       string
	Year        string
	Role        string
	Description template.HTML
	Credits     string
	Images      []string
	Links       []pubLink
	Featured    bool
	VideoURL    string
	DossierURL  string
}

type pubSection struct {
	Name        string
	Slug        string
	Description template.HTML
	CoverURL    string
	Works       []pubWork
}

type pubSite struct {
	Locale       string        // "es" | "en"
	OtherLocale  string        // el idioma inverso, para el link del toggle
	OtherURL     string        // ruta equivalente en el otro idioma
	SiteName        string
	Tagline         string
	MetaDescription string
	OGImageURL      string
	Bio             template.HTML
	Email        string
	Phone        string
	Address      string
	ShowAddress  bool
	Instagram    string
	Facebook     string
	HeroURL      string
	HeroImages   []string
	Sections     []pubSection
	Press        []pubPress
	T            map[string]string // traducciones de la UI chrome
}

type pubPress struct {
	Title, Publication, URL, Date string
	Excerpt                       template.HTML
}

func fileURL(collection, recordID, file string) string {
	if file == "" {
		return ""
	}
	return "/api/files/" + collection + "/" + recordID + "/" + file
}

// Traducciones de la UI chrome (no va a DB, hardcoded).
var uiStrings = map[string]map[string]string{
	"es": {
		"bio": "Biografía", "contact": "Contacto", "notes": "Notas",
		"address": "Dirección", "email": "Email", "phone": "Móvil",
		"networks": "Redes", "back_top": "↑ volver arriba",
		"download_dossier": "Descargar dossier (PDF)",
		"press_heading": "Notas",
		"toggle_title": "English version",
	},
	"en": {
		"bio": "Biography", "contact": "Contact", "notes": "Press",
		"address": "Address", "email": "Email", "phone": "Phone",
		"networks": "Social", "back_top": "↑ back to top",
		"download_dossier": "Download dossier (PDF)",
		"press_heading": "Press",
		"toggle_title": "Versión en español",
	},
}

// tr devuelve el valor EN si existe y el locale es "en", sino el ES.
func tr(rec interface{ GetString(string) string }, locale, field string) string {
	if locale == "en" {
		if v := rec.GetString(field + "_en"); v != "" {
			return v
		}
	}
	return rec.GetString(field)
}

func loadPubSite(app *pocketbase.PocketBase, locale string) (*pubSite, error) {
	if locale != "es" && locale != "en" {
		locale = "es"
	}
	other := "en"
	if locale == "en" {
		other = "es"
	}
	site := &pubSite{
		Locale: locale, OtherLocale: other,
		SiteName: "Amelia Repetto", Tagline: "porfolio artístico",
		T: uiStrings[locale],
	}
	if locale == "en" {
		site.OtherURL = "/"
	} else {
		site.OtherURL = "/en/"
	}

	// Settings
	settings, err := app.FindFirstRecordByFilter("site_settings", "id != ''")
	if err == nil && settings != nil {
		site.SiteName = settings.GetString("site_name")
		if t := tr(settings, locale, "tagline"); t != "" {
			site.Tagline = t
		}
		if md := tr(settings, locale, "meta_description"); md != "" {
			site.MetaDescription = md
		} else {
			site.MetaDescription = site.SiteName + " · " + site.Tagline
		}
		bioField := "bio_es"
		if locale == "en" {
			if v := settings.GetString("bio_en"); v != "" {
				bioField = "bio_en"
			}
		}
		site.Bio = safeHTML(settings.GetString(bioField))
		site.Email = settings.GetString("email")
		site.Phone = settings.GetString("phone")
		site.Address = settings.GetString("address")
		site.ShowAddress = settings.GetBool("show_address")
		site.Instagram = settings.GetString("instagram")
		site.Facebook = settings.GetString("facebook")
		if hero := settings.GetString("hero_image"); hero != "" {
			site.HeroURL = fileURL("site_settings", settings.Id, hero)
		} else if prof := settings.GetString("profile_image"); prof != "" {
			site.HeroURL = fileURL("site_settings", settings.Id, prof)
		}
		for _, f := range settings.GetStringSlice("hero_images") {
			site.HeroImages = append(site.HeroImages, fileURL("site_settings", settings.Id, f))
		}
		if og := settings.GetString("og_image"); og != "" {
			site.OGImageURL = fileURL("site_settings", settings.Id, og)
		} else if len(site.HeroImages) > 0 {
			site.OGImageURL = site.HeroImages[0]
		} else if site.HeroURL != "" {
			site.OGImageURL = site.HeroURL
		}
	}

	// Sections
	sections, err := app.FindRecordsByFilter("sections", "active = true", "+sort_order", 100, 0)
	if err != nil {
		return site, err
	}

	for _, sec := range sections {
		ps := pubSection{
			Name:        tr(sec, locale, "name"),
			Slug:        sec.GetString("slug"),
			Description: safeHTML(tr(sec, locale, "description")),
		}
		if cov := sec.GetString("cover_image"); cov != "" {
			ps.CoverURL = fileURL("sections", sec.Id, cov)
		}
		works, err := app.FindRecordsByFilter("works",
			"section = {:sec} && active = true",
			"+sort_order", 500, 0,
			map[string]any{"sec": sec.Id})
		if err == nil {
			for _, w := range works {
				pw := pubWork{
					Title:       tr(w, locale, "title"),
					Year:        w.GetString("year"),
					Role:        tr(w, locale, "role"),
					Description: safeHTML(tr(w, locale, "description")),
					Credits:     tr(w, locale, "credits"),
					Featured:    w.GetBool("featured"),
				}
				for _, f := range w.GetStringSlice("images") {
					pw.Images = append(pw.Images, fileURL("works", w.Id, f))
				}
				if v := w.GetString("video"); v != "" {
					pw.VideoURL = fileURL("works", w.Id, v)
				}
				if d := w.GetString("dossier"); d != "" {
					pw.DossierURL = fileURL("works", w.Id, d)
				}
				links, _ := app.FindRecordsByFilter("work_links",
					"work = {:w}", "+sort_order", 100, 0,
					map[string]any{"w": w.Id})
				for _, l := range links {
					pw.Links = append(pw.Links, pubLink{
						Label: l.GetString("label"),
						URL:   l.GetString("url"),
						Kind:  l.GetString("kind"),
					})
				}
				ps.Works = append(ps.Works, pw)
			}
		}
		// Ordenar por año más reciente (desc), manteniendo sort_order como desempate.
		sort.SliceStable(ps.Works, func(i, j int) bool {
			return latestYear(ps.Works[i].Year) > latestYear(ps.Works[j].Year)
		})
		site.Sections = append(site.Sections, ps)
	}

	// Press
	press, err := app.FindRecordsByFilter("press", "active = true", "+sort_order", 100, 0)
	if err == nil {
		for _, p := range press {
			site.Press = append(site.Press, pubPress{
				Title:       tr(p, locale, "title"),
				Publication: p.GetString("publication"),
				URL:         p.GetString("url"),
				Date:        p.GetString("date"),
				Excerpt:     safeHTML(tr(p, locale, "excerpt")),
			})
		}
	}

	// Ordenar secciones por orden canónico si lo preferís; ya vienen sorteadas.
	sort.SliceStable(site.Sections, func(i, j int) bool { return false })

	return site, nil
}

var publicTmpl = template.Must(template.New("public").Funcs(template.FuncMap{
	"upper": strings.ToUpper,
}).Parse(publicHTML))

func publicRobots(e *core.RequestEvent) error {
	host := e.Request.Host
	if fwd := e.Request.Header.Get("X-Forwarded-Host"); fwd != "" {
		host = fwd
	}
	scheme := "https"
	if strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "localhost") {
		scheme = "http"
	}
	body := "# robots.txt — portfolio de Amelia Repetto\n" +
		"User-agent: *\n" +
		"Allow: /\n" +
		"Allow: /en/\n" +
		"Disallow: /admin\n" +
		"Disallow: /admin/\n" +
		"Disallow: /_/\n" +
		"Disallow: /api/\n" +
		"\n" +
		"Sitemap: " + scheme + "://" + host + "/sitemap.xml\n"
	e.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	e.Response.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = e.Response.Write([]byte(body))
	return nil
}

func publicHome(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	locale := "es"
	p := e.Request.URL.Path
	if p == "/en" || p == "/en/" || strings.HasPrefix(p, "/en/") {
		locale = "en"
	}
	site, err := loadPubSite(app, locale)
	if err != nil {
		return err
	}
	// URL absoluta para og:image / twitter:image.
	if site.OGImageURL != "" && strings.HasPrefix(site.OGImageURL, "/") {
		host := e.Request.Host
		if fwd := e.Request.Header.Get("X-Forwarded-Host"); fwd != "" {
			host = fwd
		}
		// https salvo loopback de desarrollo.
		scheme := "https"
		if strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "localhost") {
			scheme = "http"
		}
		site.OGImageURL = scheme + "://" + host + site.OGImageURL
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	e.Response.WriteHeader(http.StatusOK)
	return publicTmpl.Execute(e.Response, site)
}

const publicHTML = `<!DOCTYPE html>
<html lang="{{.Locale}}">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.SiteName}} · {{.Tagline}}</title>
<meta name="description" content="{{.MetaDescription}}">
<link rel="alternate" hreflang="es" href="/">
<link rel="alternate" hreflang="en" href="/en/">
<link rel="alternate" hreflang="x-default" href="/">

<!-- Open Graph (Facebook, WhatsApp, LinkedIn) -->
<meta property="og:type" content="website">
<meta property="og:site_name" content="{{.SiteName}}">
<meta property="og:title" content="{{.SiteName}} · {{.Tagline}}">
<meta property="og:description" content="{{.MetaDescription}}">
<meta property="og:locale" content="{{if eq .Locale "es"}}es_ES{{else}}en_US{{end}}">
{{if .OGImageURL}}<meta property="og:image" content="{{.OGImageURL}}">
<meta property="og:image:alt" content="{{.SiteName}}">{{end}}

<!-- Twitter card -->
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.SiteName}} · {{.Tagline}}">
<meta name="twitter:description" content="{{.MetaDescription}}">
{{if .OGImageURL}}<meta name="twitter:image" content="{{.OGImageURL}}">{{end}}
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Roboto+Slab:wght@100..900&display=swap" rel="stylesheet">
<script src="https://cdn.tailwindcss.com"></script>
<script>
  tailwind.config = {
    theme: {
      extend: {
        fontFamily: { slab: ['"Roboto Slab"', 'Georgia', 'serif'] },
        colors: { ink: '#111111', muted: '#6a6a6a', line: '#e6e6e6' },
        screens: { nav: '1080px' },
      },
    },
  };
</script>
<style>
  *,*::before,*::after{margin:0;padding:0;box-sizing:border-box}
  :root{
    --bg:#ffffff;
    --ink:#111111;
    --muted:#6a6a6a;
    --line:#e6e6e6;
    --accent:#111111;
    --max:1100px;
  }
  html{scroll-behavior:smooth}
  body{
    background:var(--bg);color:var(--ink);
    font-family:'Roboto Slab',Georgia,serif;
    font-weight:300;line-height:1.6;
    -webkit-font-smoothing:antialiased;
  }
  a{color:inherit;text-decoration:none;border-bottom:1px solid var(--line);transition:border-color .2s}
  a:hover{border-color:var(--ink)}
  .wrap{max-width:var(--max);margin:0 auto;padding:0 1.5rem}

  /* Nav links — colores por sección (tomados del InDesign/PDF) */
  .nav-link{
    border:none !important;
    font-variation-settings:'wght' 400;
    transition:color .25s ease, font-variation-settings .25s ease, letter-spacing .25s ease;
  }
  .nav-link:hover{
    color:#6a6a6a;
    font-variation-settings:'wght' 650;
  }
  .nav-link[data-k="bio"]:hover{color:#E91E88}              /* magenta */
  .nav-link[data-k="performance"]:hover{color:#FF6B6B}      /* coral */
  .nav-link[data-k="teatro"]:hover{color:#00A676}           /* verde-teal */
  .nav-link[data-k="cine"]:hover{color:#00BCD4}             /* cian */
  .nav-link[data-k="videoclips"]:hover{color:#F2C94C}       /* amarillo oro */
  .nav-link[data-k="talleres"]:hover{color:#F5A623}         /* naranja */
  .nav-link[data-k="salud-mental"]:hover{color:#9B51E0}     /* violeta */
  .nav-link[data-k="press"]:hover{color:#C2185B}            /* rosa oscuro */
  .nav-link[data-k="contacto"]:hover{color:#111}

  /* Hero */
  .hero{padding:5rem 0 3rem}
  .hero h1{font-weight:900;font-size:clamp(2.8rem,7vw,5.5rem);line-height:.95;letter-spacing:-.02em}
  .hero .tag{font-style:italic;font-weight:300;color:var(--muted);margin-top:1rem;font-size:1.3rem}
  .hero .cover{margin-top:2.5rem;border:1px solid var(--line)}
  .hero .cover img{display:block;width:100%;height:auto;max-height:72vh;object-fit:cover}

  /* Glitch colorido sobre el hero — tamaño fijo, object-fit cover */
  .hero .cover.glitch{
    position:relative;overflow:hidden;background:#0a0a0a;
    width:100%;
    aspect-ratio: 4/5;
    max-height: 80vh;
  }
  .hero .cover.glitch img{
    position:absolute;inset:0;
    width:100% !important;height:100% !important;
    max-height:none;
    object-fit:cover;object-position:center;
    display:block;
  }
  .glitch .g-base{position:absolute}
  .glitch .g-r,.glitch .g-c{mix-blend-mode:screen;pointer-events:none;opacity:0}
  .glitch .g-r{filter:hue-rotate(-35deg) saturate(2.2) contrast(1.1)}
  .glitch .g-c{filter:hue-rotate(170deg) saturate(2.2) contrast(1.1)}
  .glitch.glitching .g-r{animation:g-red .32s linear}
  .glitch.glitching .g-c{animation:g-cyan .32s linear}
  .glitch::after{content:"";position:absolute;inset:0;pointer-events:none;background:repeating-linear-gradient(to bottom,rgba(255,255,255,0) 0 2px,rgba(0,0,0,.08) 2px 3px);opacity:0}
  .glitch.glitching::after{animation:g-scan .32s linear}
  @keyframes g-red{
    0%,100%{opacity:0;transform:translate(0,0);clip-path:inset(0 0 0 0)}
    20%{opacity:.9;transform:translate(-6px,2px);clip-path:inset(12% 0 72% 0)}
    40%{opacity:.9;transform:translate(7px,-2px);clip-path:inset(58% 0 14% 0)}
    60%{opacity:.9;transform:translate(-3px,1px);clip-path:inset(34% 0 40% 0)}
    80%{opacity:.9;transform:translate(5px,-1px);clip-path:inset(70% 0 6% 0)}
  }
  @keyframes g-cyan{
    0%,100%{opacity:0;transform:translate(0,0);clip-path:inset(0 0 0 0)}
    20%{opacity:.9;transform:translate(5px,-2px);clip-path:inset(40% 0 30% 0)}
    40%{opacity:.9;transform:translate(-6px,2px);clip-path:inset(8% 0 76% 0)}
    60%{opacity:.9;transform:translate(4px,0);clip-path:inset(55% 0 20% 0)}
    80%{opacity:.9;transform:translate(-4px,1px);clip-path:inset(24% 0 55% 0)}
  }
  @keyframes g-scan{
    0%,100%{opacity:0}
    30%,70%{opacity:.5}
  }
  @media (prefers-reduced-motion: reduce){
    .glitch .g-r,.glitch .g-c,.glitch::after{animation:none}
  }

  /* Section */
  section.page{padding:4rem 0 2rem;border-top:1px solid var(--line)}
  section.page h2{
    font-weight:900;font-size:clamp(1.8rem,4vw,2.8rem);letter-spacing:-.01em;
    margin-bottom:2rem;text-transform:uppercase;
  }
  section.page h2::after{
    content:"";display:block;width:48px;height:3px;background:var(--ink);margin-top:.6rem;
  }

  /* Work list */
  .works{display:grid;grid-template-columns:1fr;gap:2.5rem}
  @media(min-width:720px){.works{grid-template-columns:1fr 1fr}}
  .work{padding-bottom:2rem;border-bottom:1px solid var(--line)}
  .work .year{font-style:italic;color:var(--muted);font-size:.95rem;letter-spacing:.02em}
  .work h3{font-weight:700;font-size:1.4rem;line-height:1.2;margin:.3rem 0 .5rem}
  .work .role{font-size:.9rem;color:var(--muted);text-transform:uppercase;letter-spacing:.08em;margin-bottom:.7rem}
  .work .desc{font-size:.97rem;margin-bottom:.6rem}
  .work .credits{font-size:.86rem;color:var(--muted);margin-bottom:.6rem}
  .work .links{display:flex;flex-wrap:wrap;gap:.6rem 1rem;margin-top:.4rem}
  .work .links a{font-size:.82rem;letter-spacing:.06em;text-transform:uppercase}
  .work .imgs{margin-top:.9rem;display:grid;grid-template-columns:repeat(auto-fit,minmax(140px,1fr));gap:.4rem}
  .work .imgs img{width:100%;height:auto;display:block;border:1px solid var(--line);cursor:zoom-in;transition:opacity .2s}
  .work .imgs img:hover{opacity:.85}

  /* Lightbox modal con carrusel Flowbite */
  .lb-backdrop{
    position:fixed;inset:0;z-index:200;background:rgba(10,10,10,.92);
    display:none;align-items:center;justify-content:center;
  }
  .lb-backdrop.open{display:flex}
  body.lb-noscroll{overflow:hidden}
  .lb-frame{position:relative;width:100%;height:100%;max-width:1400px;padding:3rem 1rem}
  .lb-slide{position:absolute;inset:0;display:flex;align-items:center;justify-content:center;padding:3rem 1rem}
  .lb-slide img{max-width:100%;max-height:100%;object-fit:contain;display:block;border:1px solid rgba(255,255,255,.08);background:#000}
  .lb-close{
    position:absolute;top:.8rem;right:.8rem;z-index:2;
    width:42px;height:42px;border-radius:9999px;
    background:rgba(255,255,255,.1);border:1px solid rgba(255,255,255,.2);
    color:#fff;font-size:1.4rem;line-height:1;cursor:pointer;
    display:flex;align-items:center;justify-content:center;
  }
  .lb-close:hover{background:rgba(255,255,255,.2)}
  .lb-prev,.lb-next{
    position:absolute;top:50%;transform:translateY(-50%);z-index:2;
    width:48px;height:48px;border-radius:9999px;
    background:rgba(255,255,255,.1);border:1px solid rgba(255,255,255,.2);
    color:#fff;font-size:1.6rem;line-height:1;cursor:pointer;
    display:flex;align-items:center;justify-content:center;
  }
  .lb-prev:hover,.lb-next:hover{background:rgba(255,255,255,.2)}
  .lb-prev{left:1rem}.lb-next{right:1rem}
  .lb-indicators{
    position:absolute;bottom:1rem;left:50%;transform:translateX(-50%);z-index:2;
    display:flex;gap:.4rem;
  }
  .lb-indicators button{
    width:10px;height:10px;border-radius:9999px;
    background:rgba(255,255,255,.35);border:0;cursor:pointer;padding:0;
  }
  .lb-indicators button[aria-current="true"]{background:#fff}
  .lb-counter{
    position:absolute;bottom:1rem;right:1rem;z-index:2;
    color:rgba(255,255,255,.7);font-size:.8rem;letter-spacing:.05em;
  }
  @media (max-width:640px){
    .lb-prev,.lb-next{width:40px;height:40px;font-size:1.3rem}
    .lb-frame,.lb-slide{padding:3.2rem .6rem 2.6rem}
  }

  /* Bio */
  .bio{display:grid;grid-template-columns:1fr;gap:2rem;align-items:start}
  @media(min-width:820px){.bio{grid-template-columns:1fr 2fr}}
  .bio .text p{margin-bottom:.9rem;font-size:.98rem}
  .bio .text a{border-bottom:1px solid var(--line)}
  .contact p{margin-bottom:.4rem;font-size:.95rem}
  .contact .label{color:var(--muted);text-transform:uppercase;letter-spacing:.1em;font-size:.75rem;margin-top:1rem}

  /* Press */
  .press-list{display:grid;gap:1.2rem}
  .press-list article{padding:1rem 0;border-bottom:1px solid var(--line)}
  .press-list h4{font-weight:500;font-size:1.05rem;line-height:1.3}
  .press-list .meta{color:var(--muted);font-size:.82rem;text-transform:uppercase;letter-spacing:.08em;margin-top:.3rem}

  /* Footer */
  footer{padding:4rem 0 3rem;border-top:1px solid var(--line);color:var(--muted);font-size:.85rem;text-align:center;margin-top:4rem}
  footer a{border:none}
</style>
</head>
<body>

<nav class="sticky top-0 z-20 bg-white/95 backdrop-blur border-b border-line font-slab">
  <div class="max-w-[1100px] mx-auto px-6 flex items-center justify-between h-14">
    <a href="#top" class="nav-link font-bold tracking-wide text-ink">{{.SiteName}}</a>

    <!-- Desktop links -->
    <ul class="hidden nav:flex items-center gap-6 list-none">
      <li><a data-k="bio" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#bio">{{.T.bio}}</a></li>
      {{range .Sections}}<li><a data-k="{{.Slug}}" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#{{.Slug}}">{{.Name}}</a></li>{{end}}
      {{if .Press}}<li><a data-k="press" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#press">{{.T.notes}}</a></li>{{end}}
      <li><a data-k="contacto" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#contacto">{{.T.contact}}</a></li>
      <li><a href="{{.OtherURL}}" title="{{.T.toggle_title}}" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase border border-line px-2 py-1 rounded">{{if eq .Locale "es"}}EN{{else}}ES{{end}}</a></li>
    </ul>

    <!-- Mobile burger -->
    <button id="nav-toggle" aria-label="Menu" aria-expanded="false"
      class="nav:hidden inline-flex items-center justify-center w-10 h-10 -mr-2 text-ink">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M4 7h16M4 12h16M4 17h16"/></svg>
    </button>
  </div>

</nav>

<!-- Mobile full-screen overlay (outside <nav> so it can cover everything) -->
<div id="nav-panel"
  class="fixed inset-0 z-30 bg-white opacity-0 pointer-events-none transition-opacity duration-300 ease-out nav:hidden">
  <button id="nav-close" aria-label="Cerrar"
    class="absolute top-3 right-4 w-10 h-10 inline-flex items-center justify-center text-ink">
    <svg class="w-7 h-7" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 6l12 12M18 6L6 18"/></svg>
  </button>
  <ul class="h-full w-full flex flex-col items-center justify-center gap-6 list-none px-6 font-slab">
    <li><a data-k="bio" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#bio">{{.T.bio}}</a></li>
    {{range .Sections}}<li><a data-k="{{.Slug}}" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#{{.Slug}}">{{.Name}}</a></li>{{end}}
    {{if .Press}}<li><a data-k="press" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#press">{{.T.notes}}</a></li>{{end}}
    <li><a data-k="contacto" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#contacto">{{.T.contact}}</a></li>
    <li><a href="{{.OtherURL}}" class="nav-link text-xl tracking-[0.2em] uppercase text-ink border border-line px-3 py-1 rounded">{{if eq .Locale "es"}}EN{{else}}ES{{end}}</a></li>
  </ul>
</div>
<script>
  (function(){
    var btn = document.getElementById('nav-toggle');
    var closeBtn = document.getElementById('nav-close');
    var panel = document.getElementById('nav-panel');
    if(!btn || !panel) return;
    var isOpen = false;
    function close(){
      isOpen = false;
      panel.classList.remove('opacity-100','pointer-events-auto');
      panel.classList.add('opacity-0','pointer-events-none');
      btn.setAttribute('aria-expanded','false');
      document.body.style.overflow = '';
    }
    function open(){
      isOpen = true;
      panel.classList.remove('opacity-0','pointer-events-none');
      panel.classList.add('opacity-100','pointer-events-auto');
      btn.setAttribute('aria-expanded','true');
      document.body.style.overflow = 'hidden';
    }
    btn.addEventListener('click', function(){ isOpen ? close() : open(); });
    if(closeBtn) closeBtn.addEventListener('click', close);
    panel.querySelectorAll('a').forEach(function(a){ a.addEventListener('click', close); });
    document.addEventListener('keydown', function(e){ if(e.key === 'Escape' && isOpen) close(); });
  })();
</script>

<main class="wrap">

  <header class="hero" id="top">
    <h1>{{.SiteName}}</h1>
    <div class="tag">{{.Tagline}}</div>
    {{if .HeroURL}}
    <figure class="cover glitch" id="hero-glitch">
      <img class="g-base" src="{{.HeroURL}}" alt="{{.SiteName}}">
      <img class="g-r" src="{{.HeroURL}}" alt="" aria-hidden="true">
      <img class="g-c" src="{{.HeroURL}}" alt="" aria-hidden="true">
    </figure>
    <script>
      (function(){
        var el = document.getElementById('hero-glitch'); if(!el) return;
        var srcs = [{{range $i, $u := .HeroImages}}{{if $i}},{{end}}"{{$u}}"{{end}}];
        if(!srcs.length) srcs = ["{{.HeroURL}}"];
        srcs.forEach(function(s){ var i = new Image(); i.src = s; });
        var imgs = el.querySelectorAll('img');
        var idx = 0;
        function tick(){
          idx = (idx + 1) % srcs.length;
          el.classList.add('glitching');
          setTimeout(function(){
            imgs.forEach(function(im){ im.src = srcs[idx]; });
          }, 120);
          setTimeout(function(){ el.classList.remove('glitching'); }, 340);
        }
        if(srcs.length > 1) setInterval(tick, 1000);
      })();
    </script>
    {{end}}
  </header>

  <section class="page" id="bio">
    <h2>{{.T.bio}}</h2>
    <div class="bio">
      <aside class="contact" id="contacto">
        {{if .Email}}<p class="label">{{.T.email}}</p><p><a href="mailto:{{.Email}}">{{.Email}}</a></p>{{end}}
        {{if .Phone}}<p class="label">{{.T.phone}}</p><p>{{.Phone}}</p>{{end}}
        {{if and .Address .ShowAddress}}<p class="label">{{.T.address}}</p><p>{{.Address}}</p>{{end}}
        <p class="label">{{.T.networks}}</p>
        {{if .Instagram}}<p><a href="{{.Instagram}}" target="_blank" rel="noopener">Instagram</a></p>{{end}}
        {{if .Facebook}}<p><a href="{{.Facebook}}" target="_blank" rel="noopener">Facebook</a></p>{{end}}
      </aside>
      <div class="text">{{.Bio}}</div>
    </div>
  </section>

  {{range .Sections}}
  <section class="page" id="{{.Slug}}">
    <h2>{{.Name}}</h2>
    {{if .CoverURL}}<figure class="section-cover" style="margin:0 0 2rem;border:1px solid var(--line)"><img src="{{.CoverURL}}" alt="{{.Name}}" style="display:block;width:100%;height:auto;max-height:50vh;object-fit:cover"></figure>{{end}}
    {{if .Description}}<div class="section-desc" style="font-size:1.05rem;color:var(--muted);max-width:720px;margin:0 0 2.2rem">{{.Description}}</div>{{end}}
    <div class="works">
      {{range .Works}}
      <article class="work">
        {{if .Year}}<div class="year">{{.Year}}</div>{{end}}
        <h3>{{.Title}}</h3>
        {{if .Role}}<div class="role">{{.Role}}</div>{{end}}
        {{if .Description}}<div class="desc">{{.Description}}</div>{{end}}
        {{if .Credits}}<div class="credits">{{.Credits}}</div>{{end}}
        {{if .VideoURL}}
        <div class="video" style="margin:1rem 0;border:1px solid var(--line)">
          <video controls playsinline preload="metadata" style="display:block;width:100%;height:auto;max-height:70vh;background:#000">
            <source src="{{.VideoURL}}" type="video/mp4">
            <source src="{{.VideoURL}}" type="video/quicktime">
            Tu navegador no soporta el video embebido.
          </video>
        </div>
        {{end}}
        {{if .Links}}
        <div class="links">
          {{range .Links}}<a href="{{.URL}}" target="_blank" rel="noopener">{{.Label}}</a>{{end}}
          {{if .DossierURL}}<a href="{{.DossierURL}}" target="_blank" rel="noopener">{{$.T.download_dossier}}</a>{{end}}
        </div>
        {{else}}
        {{if .DossierURL}}<div class="links"><a href="{{.DossierURL}}" target="_blank" rel="noopener">{{$.T.download_dossier}}</a></div>{{end}}
        {{end}}
        {{if .Images}}
        {{$title := .Title}}
        <div class="imgs lb-gallery" data-title="{{$title}}">
          {{range $i, $u := .Images}}<img src="{{$u}}" data-i="{{$i}}" loading="lazy" alt="{{$title}}">{{end}}
        </div>
        {{end}}
      </article>
      {{end}}
    </div>
  </section>
  {{end}}

  {{if .Press}}
  <section class="page" id="press">
    <h2>{{.T.press_heading}}</h2>
    <div class="press-list">
      {{range .Press}}
      <article>
        <h4>{{if .URL}}<a href="{{.URL}}" target="_blank" rel="noopener">{{.Title}}</a>{{else}}{{.Title}}{{end}}</h4>
        <div class="meta">{{.Publication}}{{if .Date}} · {{.Date}}{{end}}</div>
        {{if .Excerpt}}<div class="desc">{{.Excerpt}}</div>{{end}}
      </article>
      {{end}}
    </div>
  </section>
  {{end}}

</main>

<footer>
  © {{.SiteName}} · <a href="#top">{{.T.back_top}}</a> · <a href="/admin">admin</a>
</footer>

<!-- Lightbox con carrusel Flowbite -->
<div id="lb" class="lb-backdrop" role="dialog" aria-modal="true" aria-label="Galería">
  <button id="lb-close" class="lb-close" aria-label="Cerrar">×</button>
  <div id="lb-carousel" class="lb-frame relative">
    <div id="lb-inner"></div>
    <div id="lb-dots" class="lb-indicators"></div>
    <div id="lb-counter" class="lb-counter"></div>
    <button id="lb-prev" type="button" class="lb-prev" aria-label="Anterior">‹</button>
    <button id="lb-next" type="button" class="lb-next" aria-label="Siguiente">›</button>
  </div>
</div>

<script>
  (function(){
    var lb      = document.getElementById('lb');
    var inner   = document.getElementById('lb-inner');
    var dots    = document.getElementById('lb-dots');
    var counter = document.getElementById('lb-counter');
    var srcs = [], idx = 0;

    function render(){
      inner.innerHTML = '';
      dots.innerHTML = '';
      srcs.forEach(function(u, i){
        var slide = document.createElement('div');
        slide.className = 'lb-slide';
        slide.style.display = (i === idx) ? 'flex' : 'none';
        var img = document.createElement('img');
        img.src = u;
        img.alt = '';
        slide.appendChild(img);
        inner.appendChild(slide);

        var dot = document.createElement('button');
        dot.type = 'button';
        dot.setAttribute('aria-label', 'Imagen ' + (i+1));
        if (i === idx) dot.setAttribute('aria-current','true');
        dot.addEventListener('click', function(){ show(i); });
        dots.appendChild(dot);
      });
      counter.textContent = (idx+1) + ' / ' + srcs.length;
    }
    function show(n){
      idx = (n + srcs.length) % srcs.length;
      Array.from(inner.children).forEach(function(s, i){
        s.style.display = (i === idx) ? 'flex' : 'none';
      });
      Array.from(dots.children).forEach(function(b, i){
        if (i === idx) b.setAttribute('aria-current','true');
        else b.removeAttribute('aria-current');
      });
      counter.textContent = (idx+1) + ' / ' + srcs.length;
    }
    function open(list, start){
      if (!list.length) return;
      srcs = list;
      idx = start || 0;
      render();
      lb.classList.add('open');
      document.body.classList.add('lb-noscroll');
    }
    function close(){
      lb.classList.remove('open');
      document.body.classList.remove('lb-noscroll');
    }

    document.querySelectorAll('.lb-gallery').forEach(function(g){
      var imgs = Array.from(g.querySelectorAll('img'));
      var list = imgs.map(function(i){ return i.src; });
      imgs.forEach(function(el, k){
        el.addEventListener('click', function(){ open(list, k); });
      });
    });
    document.getElementById('lb-close').addEventListener('click', close);
    document.getElementById('lb-prev').addEventListener('click', function(){ show(idx-1); });
    document.getElementById('lb-next').addEventListener('click', function(){ show(idx+1); });
    lb.addEventListener('click', function(e){ if (e.target === lb) close(); });
    document.addEventListener('keydown', function(e){
      if (!lb.classList.contains('open')) return;
      if (e.key === 'Escape') close();
      else if (e.key === 'ArrowLeft') show(idx-1);
      else if (e.key === 'ArrowRight') show(idx+1);
    });
  })();
</script>

</body>
</html>
`
