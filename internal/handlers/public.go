package handlers

import (
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

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
	Name  string
	Slug  string
	Works []pubWork
}

type pubSite struct {
	SiteName  string
	Tagline   string
	BioES     template.HTML
	Email       string
	Phone       string
	Address     string
	ShowAddress bool
	Instagram   string
	Facebook    string
	HeroURL    string
	HeroImages []string
	Sections   []pubSection
	Press     []pubPress
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

func loadPubSite(app *pocketbase.PocketBase) (*pubSite, error) {
	site := &pubSite{SiteName: "Amelia Repetto", Tagline: "porfolio artístico"}

	// Settings
	settings, err := app.FindFirstRecordByFilter("site_settings", "id != ''")
	if err == nil && settings != nil {
		site.SiteName = settings.GetString("site_name")
		site.Tagline = settings.GetString("tagline")
		site.BioES = template.HTML(settings.GetString("bio_es"))
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
	}

	// Sections
	sections, err := app.FindRecordsByFilter("sections", "active = true", "+sort_order", 100, 0)
	if err != nil {
		return site, err
	}

	for _, sec := range sections {
		ps := pubSection{Name: sec.GetString("name"), Slug: sec.GetString("slug")}
		works, err := app.FindRecordsByFilter("works",
			"section = {:sec} && active = true",
			"+sort_order", 500, 0,
			map[string]any{"sec": sec.Id})
		if err == nil {
			for _, w := range works {
				pw := pubWork{
					Title:       w.GetString("title"),
					Year:        w.GetString("year"),
					Role:        w.GetString("role"),
					Description: template.HTML(w.GetString("description")),
					Credits:     w.GetString("credits"),
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
		site.Sections = append(site.Sections, ps)
	}

	// Press
	press, err := app.FindRecordsByFilter("press", "active = true", "+sort_order", 100, 0)
	if err == nil {
		for _, p := range press {
			site.Press = append(site.Press, pubPress{
				Title:       p.GetString("title"),
				Publication: p.GetString("publication"),
				URL:         p.GetString("url"),
				Date:        p.GetString("date"),
				Excerpt:     template.HTML(p.GetString("excerpt")),
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

func publicHome(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	site, err := loadPubSite(app)
	if err != nil {
		return err
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	e.Response.WriteHeader(http.StatusOK)
	return publicTmpl.Execute(e.Response, site)
}

const publicHTML = `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.SiteName}} · {{.Tagline}}</title>
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
  .work .imgs img{width:100%;height:auto;display:block;border:1px solid var(--line)}

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
      <li><a data-k="bio" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#bio">Biografía</a></li>
      {{range .Sections}}<li><a data-k="{{.Slug}}" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#{{.Slug}}">{{.Name}}</a></li>{{end}}
      {{if .Press}}<li><a data-k="press" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#press">Notas</a></li>{{end}}
      <li><a data-k="contacto" class="nav-link text-[0.78rem] tracking-[0.14em] uppercase text-ink" href="#contacto">Contacto</a></li>
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
    <li><a data-k="bio" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#bio">Biografía</a></li>
    {{range .Sections}}<li><a data-k="{{.Slug}}" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#{{.Slug}}">{{.Name}}</a></li>{{end}}
    {{if .Press}}<li><a data-k="press" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#press">Notas</a></li>{{end}}
    <li><a data-k="contacto" class="nav-link text-2xl tracking-[0.12em] uppercase text-ink" href="#contacto">Contacto</a></li>
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
    <h2>Biografía</h2>
    <div class="bio">
      <aside class="contact" id="contacto">
        {{if .Email}}<p class="label">Email</p><p><a href="mailto:{{.Email}}">{{.Email}}</a></p>{{end}}
        {{if .Phone}}<p class="label">Móvil</p><p>{{.Phone}}</p>{{end}}
        {{if and .Address .ShowAddress}}<p class="label">Dirección</p><p>{{.Address}}</p>{{end}}
        <p class="label">Redes</p>
        {{if .Instagram}}<p><a href="{{.Instagram}}" target="_blank" rel="noopener">Instagram</a></p>{{end}}
        {{if .Facebook}}<p><a href="{{.Facebook}}" target="_blank" rel="noopener">Facebook</a></p>{{end}}
      </aside>
      <div class="text">{{.BioES}}</div>
    </div>
  </section>

  {{range .Sections}}
  <section class="page" id="{{.Slug}}">
    <h2>{{.Name}}</h2>
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
          {{if .DossierURL}}<a href="{{.DossierURL}}" target="_blank" rel="noopener">Descargar dossier (PDF)</a>{{end}}
        </div>
        {{else}}
        {{if .DossierURL}}<div class="links"><a href="{{.DossierURL}}" target="_blank" rel="noopener">Descargar dossier (PDF)</a></div>{{end}}
        {{end}}
        {{if .Images}}
        <div class="imgs">
          {{range .Images}}<img src="{{.}}" loading="lazy" alt="">{{end}}
        </div>
        {{end}}
      </article>
      {{end}}
    </div>
  </section>
  {{end}}

  {{if .Press}}
  <section class="page" id="press">
    <h2>Notas / Press</h2>
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
  © {{.SiteName}} · <a href="#top">↑ volver arriba</a> · <a href="/admin">admin</a>
</footer>

</body>
</html>
`
