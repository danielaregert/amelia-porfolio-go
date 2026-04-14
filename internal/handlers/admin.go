package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"

	"porfolio-amelia/internal/adminsession"
	"porfolio-amelia/internal/ratelimit"
	"porfolio-amelia/internal/sanitize"
	"porfolio-amelia/internal/views"
)

var loginLimiter = ratelimit.NewLoginLimiter()

func requireAdmin(e *core.RequestEvent, app *pocketbase.PocketBase) (*core.Record, bool) {
	user := adminsession.CurrentUser(app, e.Request)
	if user == nil {
		http.Redirect(e.Response, e.Request, "/admin/login", http.StatusSeeOther)
		return nil, false
	}
	return user, true
}

// ==================== AUTH ====================

func adminLoginView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if u := adminsession.CurrentUser(app, e.Request); u != nil {
		http.Redirect(e.Response, e.Request, "/admin", http.StatusSeeOther)
		return nil
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = e.Response.Write([]byte(renderLogin("")))
	return nil
}

func renderLogin(errMsg string) string {
	errBlock := ""
	if errMsg != "" {
		errBlock = `<div class="err">` + errMsg + `</div>`
	}
	return `<!DOCTYPE html>
<html lang="es"><head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Admin · Amelia Repetto</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Roboto+Slab:wght@300;400;500;700;900&display=swap" rel="stylesheet">
<style>
  *,*::before,*::after{margin:0;padding:0;box-sizing:border-box}
  html,body{height:100%}
  body{
    font-family:'Roboto Slab',Georgia,serif;
    background:#fafafa;color:#111;
    display:flex;align-items:center;justify-content:center;
    min-height:100vh;min-height:100dvh;padding:1.5rem;
    -webkit-font-smoothing:antialiased;
  }
  .card{
    width:100%;max-width:380px;
    background:#fff;border:1px solid #e6e6e6;
    padding:2.8rem 2.2rem 2.2rem;
    box-shadow:0 20px 60px -20px rgba(0,0,0,.15);
  }
  .face{
    font-family:ui-monospace,SFMono-Regular,Menlo,monospace;
    text-align:center;font-size:1.9rem;letter-spacing:.04em;
    margin-bottom:.6rem;color:#111;
  }
  h1{font-weight:900;font-size:1.5rem;text-align:center;letter-spacing:-.01em;margin-bottom:.3rem}
  .sub{text-align:center;color:#6a6a6a;font-style:italic;font-size:.95rem;margin-bottom:1.8rem}
  form{display:flex;flex-direction:column;gap:1rem}
  label{font-size:.78rem;text-transform:uppercase;letter-spacing:.1em;color:#6a6a6a;display:block;margin-bottom:.35rem}
  input{
    width:100%;font-family:inherit;font-size:1rem;
    padding:.7rem .8rem;border:1px solid #d9d9d9;background:#fff;
    outline:none;transition:border-color .15s;
  }
  input:focus{border-color:#111}
  button{
    font-family:inherit;font-weight:700;font-size:.95rem;letter-spacing:.05em;
    background:#111;color:#fff;border:none;padding:.85rem 1rem;
    cursor:pointer;margin-top:.3rem;transition:background .15s;
  }
  button:hover{background:#333}
  .err{
    background:#fff1f1;border:1px solid #f3c2c2;color:#9b1c1c;
    padding:.6rem .8rem;font-size:.88rem;margin-bottom:1rem;
  }
  .foot{text-align:center;margin-top:1.5rem;font-size:.8rem;color:#9a9a9a}
  .foot a{color:#6a6a6a;border-bottom:1px solid #e6e6e6;text-decoration:none}
  .foot a:hover{color:#111}
</style>
</head><body>
  <main class="card">
    <div class="face">( ◠‿◠ )</div>
    <h1>Amelia Repetto</h1>
    <div class="sub">panel de administración</div>
    ` + errBlock + `
    <form method="POST" action="/admin/login" autocomplete="on">
      <div>
        <label for="email">Email</label>
        <input id="email" type="email" name="email" required autofocus>
      </div>
      <div>
        <label for="password">Contraseña</label>
        <input id="password" type="password" name="password" required>
      </div>
      <button type="submit">Entrar</button>
    </form>
    <div class="foot"><a href="/">← volver al porfolio</a></div>
  </main>
</body></html>`
}

func adminLoginSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	ip := ratelimit.ClientIP(e.Request)
	if allowed, retryAfter := loginLimiter.Allowed(ip); !allowed {
		mins := int(retryAfter.Minutes()) + 1
		e.Response.WriteHeader(http.StatusTooManyRequests)
		e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
		msg := fmt.Sprintf("Demasiados intentos. Probá en %d minutos.", mins)
		_, _ = e.Response.Write([]byte(renderLogin(msg)))
		return nil
	}
	if err := e.Request.ParseForm(); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	email := strings.TrimSpace(e.Request.PostFormValue("email"))
	password := e.Request.PostFormValue("password")
	_, token, err := adminsession.Authenticate(app, email, password)
	if err != nil {
		loginLimiter.RecordFailure(ip)
		e.Response.WriteHeader(http.StatusUnauthorized)
		e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = e.Response.Write([]byte(renderLogin("Email o contraseña incorrectos.")))
		return nil
	}
	loginLimiter.RecordSuccess(ip)
	adminsession.SetCookie(e.Response, e.Request, token)
	http.Redirect(e.Response, e.Request, "/admin", http.StatusSeeOther)
	return nil
}

func adminLogout(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	adminsession.ClearCookie(e.Response)
	http.Redirect(e.Response, e.Request, "/admin/login", http.StatusSeeOther)
	return nil
}

// ==================== SECTIONS ====================

func adminSectionsList(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	records, err := app.FindRecordsByFilter("sections", "", "sort_order", 100, 0)
	if err != nil {
		return e.InternalServerError("error cargando secciones", err)
	}
	rows := make([]views.AdminSectionRow, 0, len(records))
	for _, r := range records {
		works, _ := app.FindRecordsByFilter("works", "section = {:s}", "", 500, 0, map[string]any{"s": r.Id})
		rows = append(rows, views.AdminSectionRow{
			ID:         r.Id,
			Name:       r.GetString("name"),
			Slug:       r.GetString("slug"),
			SortOrder:  r.GetInt("sort_order"),
			Active:     r.GetBool("active"),
			WorksCount: len(works),
		})
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminSectionsPage(user.GetString("email"), rows).Render(context.Background(), e.Response)
}

func adminSectionNewView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	f := views.AdminSectionForm{IsNew: true, Active: true}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminSectionFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func adminSectionNewSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	if err := e.Request.ParseMultipartForm(20 * 1024 * 1024); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	f := readSectionForm(e.Request, true)
	if f.Name == "" || f.Slug == "" {
		return renderSectionForm(e, user, f, "Nombre y slug son obligatorios.")
	}
	if dup, _ := app.FindFirstRecordByFilter("sections", "slug = {:s}", map[string]any{"s": f.Slug}); dup != nil {
		return renderSectionForm(e, user, f, "Ya existe una sección con ese slug.")
	}
	col, err := app.FindCollectionByNameOrId("sections")
	if err != nil {
		return e.InternalServerError("sections collection", err)
	}
	rec := core.NewRecord(col)
	rec.Set("name", f.Name)
	rec.Set("slug", f.Slug)
	rec.Set("description", sanitize.HTML(f.Description))
	rec.Set("sort_order", f.SortOrder)
	rec.Set("active", f.Active)
	if files := uploadedFile(e.Request, "cover_image"); files != nil {
		rec.Set("cover_image", files)
	}
	if err := app.Save(rec); err != nil {
		return renderSectionForm(e, user, f, "Error guardando: "+err.Error())
	}
	http.Redirect(e.Response, e.Request, "/admin", http.StatusSeeOther)
	return nil
}

func adminSectionEditView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("sections", id)
	if err != nil {
		return e.NotFoundError("sección no encontrada", nil)
	}
	f := views.AdminSectionForm{
		IsNew:       false,
		ID:          rec.Id,
		Name:        rec.GetString("name"),
		Slug:        rec.GetString("slug"),
		Description: rec.GetString("description"),
		SortOrder:   rec.GetInt("sort_order"),
		Active:      rec.GetBool("active"),
	}
	if cover := rec.GetString("cover_image"); cover != "" {
		f.CoverURL = "/api/files/sections/" + rec.Id + "/" + cover
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminSectionFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func adminSectionEditSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("sections", id)
	if err != nil {
		return e.NotFoundError("sección no encontrada", nil)
	}
	if err := e.Request.ParseMultipartForm(20 * 1024 * 1024); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	f := readSectionForm(e.Request, false)
	f.ID = rec.Id
	if f.Name == "" || f.Slug == "" {
		return renderSectionForm(e, user, f, "Nombre y slug son obligatorios.")
	}
	if f.Slug != rec.GetString("slug") {
		if dup, _ := app.FindFirstRecordByFilter("sections", "slug = {:s}", map[string]any{"s": f.Slug}); dup != nil && dup.Id != rec.Id {
			return renderSectionForm(e, user, f, "Ya existe otra sección con ese slug.")
		}
	}
	rec.Set("name", f.Name)
	rec.Set("slug", f.Slug)
	rec.Set("description", sanitize.HTML(f.Description))
	rec.Set("sort_order", f.SortOrder)
	rec.Set("active", f.Active)
	if files := uploadedFile(e.Request, "cover_image"); files != nil {
		rec.Set("cover_image", files)
	}
	if err := app.Save(rec); err != nil {
		return renderSectionForm(e, user, f, "Error guardando: "+err.Error())
	}
	http.Redirect(e.Response, e.Request, "/admin", http.StatusSeeOther)
	return nil
}

func adminSectionToggle(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("sections", id)
	if err != nil {
		return e.NotFoundError("sección no encontrada", nil)
	}
	rec.Set("active", !rec.GetBool("active"))
	if err := app.Save(rec); err != nil {
		return e.InternalServerError("error", err)
	}
	http.Redirect(e.Response, e.Request, "/admin", http.StatusSeeOther)
	return nil
}

func adminSectionDelete(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("sections", id)
	if err != nil {
		return e.NotFoundError("sección no encontrada", nil)
	}
	if err := app.Delete(rec); err != nil {
		return e.InternalServerError("error borrando sección", err)
	}
	http.Redirect(e.Response, e.Request, "/admin", http.StatusSeeOther)
	return nil
}

// ==================== WORKS ====================

func loadSectionOptions(app *pocketbase.PocketBase) []views.SectionOption {
	records, _ := app.FindRecordsByFilter("sections", "", "sort_order", 100, 0)
	opts := make([]views.SectionOption, 0, len(records))
	for _, r := range records {
		opts = append(opts, views.SectionOption{ID: r.Id, Name: r.GetString("name")})
	}
	return opts
}

func adminWorksList(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	sectionFilter := e.Request.URL.Query().Get("section")
	filter := ""
	params := map[string]any{}
	if sectionFilter != "" {
		filter = "section = {:s}"
		params["s"] = sectionFilter
	}
	records, err := app.FindRecordsByFilter("works", filter, "-year, sort_order", 500, 0, params)
	if err != nil {
		return e.InternalServerError("error cargando obras", err)
	}

	// Cache section names
	sectionNames := map[string]string{}
	sections, _ := app.FindRecordsByFilter("sections", "", "sort_order", 100, 0)
	for _, s := range sections {
		sectionNames[s.Id] = s.GetString("name")
	}

	rows := make([]views.AdminWorkRow, 0, len(records))
	for _, r := range records {
		thumb := ""
		if imgs := r.GetStringSlice("images"); len(imgs) > 0 {
			thumb = "/api/files/works/" + r.Id + "/" + imgs[0]
		}
		links, _ := app.FindRecordsByFilter("work_links", "work = {:w}", "", 50, 0, map[string]any{"w": r.Id})
		sID := r.GetString("section")
		rows = append(rows, views.AdminWorkRow{
			ID:          r.Id,
			Title:       r.GetString("title"),
			Slug:        r.GetString("slug"),
			Year:        r.GetString("year"),
			Role:        r.GetString("role"),
			SectionName: sectionNames[sID],
			SectionID:   sID,
			ThumbURL:    thumb,
			Active:      r.GetBool("active"),
			Featured:    r.GetBool("featured"),
			LinksCount:  len(links),
		})
	}

	sectionOpts := loadSectionOptions(app)
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminWorksPage(user.GetString("email"), rows, sectionFilter, sectionOpts).Render(context.Background(), e.Response)
}

func adminWorkNewView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	f := views.AdminWorkForm{
		IsNew:    true,
		Active:   true,
		Sections: loadSectionOptions(app),
	}
	// Pre-select section if coming from section filter
	if sid := e.Request.URL.Query().Get("section"); sid != "" {
		f.SectionID = sid
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminWorkFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func adminWorkNewSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	if err := e.Request.ParseMultipartForm(100 * 1024 * 1024); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	f := readWorkForm(e.Request, true)
	f.Sections = loadSectionOptions(app)
	if f.Title == "" || f.Slug == "" || f.SectionID == "" {
		return renderWorkForm(e, user, f, "Título, slug y sección son obligatorios.")
	}
	if dup, _ := app.FindFirstRecordByFilter("works", "slug = {:s}", map[string]any{"s": f.Slug}); dup != nil {
		return renderWorkForm(e, user, f, "Ya existe una obra con ese slug.")
	}
	col, err := app.FindCollectionByNameOrId("works")
	if err != nil {
		return e.InternalServerError("works collection", err)
	}
	rec := core.NewRecord(col)
	rec.Set("section", f.SectionID)
	rec.Set("title", f.Title)
	rec.Set("slug", f.Slug)
	rec.Set("year", f.Year)
	rec.Set("role", f.Role)
	rec.Set("description", sanitize.HTML(f.Description))
	rec.Set("credits", f.Credits)
	rec.Set("sort_order", f.SortOrder)
	rec.Set("active", f.Active)
	rec.Set("featured", f.Featured)
	if files := uploadedFiles(e.Request, "images"); len(files) > 0 {
		rec.Set("images", files)
	}
	if err := app.Save(rec); err != nil {
		return renderWorkForm(e, user, f, "Error guardando: "+err.Error())
	}
	http.Redirect(e.Response, e.Request, "/admin/works/"+rec.Id, http.StatusSeeOther)
	return nil
}

func adminWorkEditView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("works", id)
	if err != nil {
		return e.NotFoundError("obra no encontrada", nil)
	}
	f := views.AdminWorkForm{
		IsNew:       false,
		ID:          rec.Id,
		SectionID:   rec.GetString("section"),
		Title:       rec.GetString("title"),
		Slug:        rec.GetString("slug"),
		Year:        rec.GetString("year"),
		Role:        rec.GetString("role"),
		Description: rec.GetString("description"),
		Credits:     rec.GetString("credits"),
		SortOrder:   rec.GetInt("sort_order"),
		Active:      rec.GetBool("active"),
		Featured:    rec.GetBool("featured"),
		Sections:    loadSectionOptions(app),
	}
	for _, name := range rec.GetStringSlice("images") {
		f.Images = append(f.Images, views.AdminImage{
			URL:      "/api/files/works/" + rec.Id + "/" + name,
			Filename: name,
		})
	}
	// Load links
	links, _ := app.FindRecordsByFilter("work_links", "work = {:w}", "sort_order", 50, 0, map[string]any{"w": rec.Id})
	for _, l := range links {
		f.Links = append(f.Links, views.AdminWorkLink{
			ID:    l.Id,
			Label: l.GetString("label"),
			URL:   l.GetString("url"),
			Kind:  l.GetString("kind"),
		})
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminWorkFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func adminWorkEditSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("works", id)
	if err != nil {
		return e.NotFoundError("obra no encontrada", nil)
	}
	if err := e.Request.ParseMultipartForm(100 * 1024 * 1024); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	f := readWorkForm(e.Request, false)
	f.ID = rec.Id
	f.Sections = loadSectionOptions(app)
	if f.Title == "" || f.Slug == "" || f.SectionID == "" {
		return renderWorkForm(e, user, f, "Título, slug y sección son obligatorios.")
	}
	if f.Slug != rec.GetString("slug") {
		if dup, _ := app.FindFirstRecordByFilter("works", "slug = {:s}", map[string]any{"s": f.Slug}); dup != nil && dup.Id != rec.Id {
			return renderWorkForm(e, user, f, "Ya existe otra obra con ese slug.")
		}
	}
	rec.Set("section", f.SectionID)
	rec.Set("title", f.Title)
	rec.Set("slug", f.Slug)
	rec.Set("year", f.Year)
	rec.Set("role", f.Role)
	rec.Set("description", sanitize.HTML(f.Description))
	rec.Set("credits", f.Credits)
	rec.Set("sort_order", f.SortOrder)
	rec.Set("active", f.Active)
	rec.Set("featured", f.Featured)
	if files := uploadedFiles(e.Request, "images"); len(files) > 0 {
		rec.Set("+images", files)
	}
	if file := uploadedFile(e.Request, "video"); file != nil {
		rec.Set("video", file)
	} else if e.Request.PostFormValue("video_delete") == "1" {
		rec.Set("video", nil)
	}
	if file := uploadedFile(e.Request, "dossier"); file != nil {
		rec.Set("dossier", file)
	} else if e.Request.PostFormValue("dossier_delete") == "1" {
		rec.Set("dossier", nil)
	}
	if err := app.Save(rec); err != nil {
		return renderWorkForm(e, user, f, "Error guardando: "+err.Error())
	}
	http.Redirect(e.Response, e.Request, "/admin/works/"+rec.Id, http.StatusSeeOther)
	return nil
}

func adminWorkDelete(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("works", id)
	if err != nil {
		return e.NotFoundError("obra no encontrada", nil)
	}
	if err := app.Delete(rec); err != nil {
		return e.InternalServerError("error borrando", err)
	}
	http.Redirect(e.Response, e.Request, "/admin/works", http.StatusSeeOther)
	return nil
}

func adminWorkImageDelete(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("works", id)
	if err != nil {
		return e.NotFoundError("obra no encontrada", nil)
	}
	if err := e.Request.ParseForm(); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	filename := e.Request.PostFormValue("filename")
	if filename == "" {
		return e.BadRequestError("filename requerido", nil)
	}
	rec.Set("images-", filename)
	if err := app.Save(rec); err != nil {
		return e.InternalServerError("error borrando imagen", err)
	}
	var imgs []views.AdminImage
	for _, n := range rec.GetStringSlice("images") {
		imgs = append(imgs, views.AdminImage{
			URL:      "/api/files/works/" + rec.Id + "/" + n,
			Filename: n,
		})
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminWorkImagesGrid(rec.Id, imgs).Render(context.Background(), e.Response)
}

// ==================== WORK LINKS ====================

func adminWorkLinkNewForm(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminWorkLinkForm(id).Render(context.Background(), e.Response)
}

func adminWorkLinkCreate(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	workID := e.Request.PathValue("id")
	if err := e.Request.ParseForm(); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	label := strings.TrimSpace(e.Request.PostFormValue("label"))
	url := strings.TrimSpace(e.Request.PostFormValue("url"))
	kind := e.Request.PostFormValue("kind")
	if label == "" || url == "" {
		return e.BadRequestError("label y url requeridos", nil)
	}
	col, err := app.FindCollectionByNameOrId("work_links")
	if err != nil {
		return e.InternalServerError("work_links collection", err)
	}
	rec := core.NewRecord(col)
	rec.Set("work", workID)
	rec.Set("label", label)
	rec.Set("url", url)
	rec.Set("kind", kind)
	if err := app.Save(rec); err != nil {
		return e.InternalServerError("error guardando link", err)
	}
	l := views.AdminWorkLink{ID: rec.Id, Label: label, URL: url, Kind: kind}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminWorkLinkRow(workID, l).Render(context.Background(), e.Response)
}

func adminWorkLinkDelete(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	lid := e.Request.PathValue("lid")
	rec, err := app.FindRecordById("work_links", lid)
	if err != nil {
		return e.NotFoundError("link no encontrado", nil)
	}
	if err := app.Delete(rec); err != nil {
		return e.InternalServerError("error borrando link", err)
	}
	// Return empty to remove the element
	return nil
}

// ==================== PRESS ====================

func adminPressList(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	records, err := app.FindRecordsByFilter("press", "", "sort_order", 500, 0)
	if err != nil {
		return e.InternalServerError("error cargando prensa", err)
	}
	rows := make([]views.AdminPressRow, 0, len(records))
	for _, r := range records {
		rows = append(rows, views.AdminPressRow{
			ID:          r.Id,
			Title:       r.GetString("title"),
			Publication: r.GetString("publication"),
			URL:         r.GetString("url"),
			Date:        r.GetString("date"),
			Active:      r.GetBool("active"),
		})
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminPressPage(user.GetString("email"), rows).Render(context.Background(), e.Response)
}

func adminPressNewView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	f := views.AdminPressForm{IsNew: true, Active: true}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminPressFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func adminPressNewSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	if err := e.Request.ParseForm(); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	f := readPressForm(e.Request, true)
	if f.Title == "" {
		return renderPressForm(e, user, f, "El título es obligatorio.")
	}
	col, err := app.FindCollectionByNameOrId("press")
	if err != nil {
		return e.InternalServerError("press collection", err)
	}
	rec := core.NewRecord(col)
	rec.Set("title", f.Title)
	rec.Set("publication", f.Publication)
	rec.Set("url", f.URL)
	rec.Set("date", f.Date)
	rec.Set("excerpt", sanitize.HTML(f.Excerpt))
	rec.Set("sort_order", f.SortOrder)
	rec.Set("active", f.Active)
	if err := app.Save(rec); err != nil {
		return renderPressForm(e, user, f, "Error guardando: "+err.Error())
	}
	http.Redirect(e.Response, e.Request, "/admin/press", http.StatusSeeOther)
	return nil
}

func adminPressEditView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("press", id)
	if err != nil {
		return e.NotFoundError("nota no encontrada", nil)
	}
	f := views.AdminPressForm{
		IsNew:       false,
		ID:          rec.Id,
		Title:       rec.GetString("title"),
		Publication: rec.GetString("publication"),
		URL:         rec.GetString("url"),
		Date:        rec.GetString("date"),
		Excerpt:     rec.GetString("excerpt"),
		SortOrder:   rec.GetInt("sort_order"),
		Active:      rec.GetBool("active"),
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminPressFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func adminPressEditSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("press", id)
	if err != nil {
		return e.NotFoundError("nota no encontrada", nil)
	}
	if err := e.Request.ParseForm(); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	f := readPressForm(e.Request, false)
	f.ID = rec.Id
	if f.Title == "" {
		return renderPressForm(e, user, f, "El título es obligatorio.")
	}
	rec.Set("title", f.Title)
	rec.Set("publication", f.Publication)
	rec.Set("url", f.URL)
	rec.Set("date", f.Date)
	rec.Set("excerpt", sanitize.HTML(f.Excerpt))
	rec.Set("sort_order", f.SortOrder)
	rec.Set("active", f.Active)
	if err := app.Save(rec); err != nil {
		return renderPressForm(e, user, f, "Error guardando: "+err.Error())
	}
	http.Redirect(e.Response, e.Request, "/admin/press", http.StatusSeeOther)
	return nil
}

func adminPressDelete(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	id := e.Request.PathValue("id")
	rec, err := app.FindRecordById("press", id)
	if err != nil {
		return e.NotFoundError("nota no encontrada", nil)
	}
	if err := app.Delete(rec); err != nil {
		return e.InternalServerError("error borrando", err)
	}
	http.Redirect(e.Response, e.Request, "/admin/press", http.StatusSeeOther)
	return nil
}

// ==================== SETTINGS — HERO CAROUSEL HTMX ====================

func adminHeroImageDelete(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	rec, err := getOrCreateSettings(app)
	if err != nil {
		return e.InternalServerError("error cargando config", err)
	}
	if err := e.Request.ParseForm(); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	filename := e.Request.PostFormValue("filename")
	if filename == "" {
		return e.BadRequestError("filename requerido", nil)
	}
	rec.Set("hero_images-", filename)
	if err := app.Save(rec); err != nil {
		return e.InternalServerError("error borrando imagen", err)
	}
	return writeHeroGrid(e, rec)
}

func adminHeroImageMove(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	if _, ok := requireAdmin(e, app); !ok {
		return nil
	}
	rec, err := getOrCreateSettings(app)
	if err != nil {
		return e.InternalServerError("error cargando config", err)
	}
	if err := e.Request.ParseForm(); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	filename := e.Request.PostFormValue("filename")
	dir := e.Request.PostFormValue("dir")
	if filename == "" || (dir != "up" && dir != "down") {
		return e.BadRequestError("filename y dir (up|down) requeridos", nil)
	}
	files := append([]string{}, rec.GetStringSlice("hero_images")...)
	idx := -1
	for i, f := range files {
		if f == filename {
			idx = i
			break
		}
	}
	if idx < 0 {
		return e.BadRequestError("filename no encontrado", nil)
	}
	newIdx := idx
	if dir == "up" && idx > 0 {
		newIdx = idx - 1
	} else if dir == "down" && idx < len(files)-1 {
		newIdx = idx + 1
	}
	if newIdx != idx {
		files[idx], files[newIdx] = files[newIdx], files[idx]
		rec.Set("hero_images", files)
		if err := app.Save(rec); err != nil {
			return e.InternalServerError("error reordenando", err)
		}
	}
	return writeHeroGrid(e, rec)
}

func writeHeroGrid(e *core.RequestEvent, rec *core.Record) error {
	var imgs []views.AdminImage
	for _, n := range rec.GetStringSlice("hero_images") {
		imgs = append(imgs, views.AdminImage{
			URL:      "/api/files/site_settings/" + rec.Id + "/" + n,
			Filename: n,
		})
	}
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = e.Response.Write([]byte(views.RenderHeroImagesGrid(imgs)))
	return nil
}

// ==================== SETTINGS ====================

func getOrCreateSettings(app *pocketbase.PocketBase) (*core.Record, error) {
	records, err := app.FindRecordsByFilter("site_settings", "", "", 1, 0)
	if err == nil && len(records) > 0 {
		return records[0], nil
	}
	col, err := app.FindCollectionByNameOrId("site_settings")
	if err != nil {
		return nil, err
	}
	rec := core.NewRecord(col)
	rec.Set("site_name", "Amelia Repetto")
	if err := app.Save(rec); err != nil {
		return nil, err
	}
	return rec, nil
}

func adminSettingsView(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	rec, err := getOrCreateSettings(app)
	if err != nil {
		return e.InternalServerError("error cargando config", err)
	}
	f := settingsToForm(rec)
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminSettingsPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func adminSettingsSubmit(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	user, ok := requireAdmin(e, app)
	if !ok {
		return nil
	}
	rec, err := getOrCreateSettings(app)
	if err != nil {
		return e.InternalServerError("error cargando config", err)
	}
	if err := e.Request.ParseMultipartForm(100 * 1024 * 1024); err != nil {
		return e.BadRequestError("form inválido", err)
	}
	rec.Set("site_name", strings.TrimSpace(e.Request.PostFormValue("site_name")))
	rec.Set("tagline", strings.TrimSpace(e.Request.PostFormValue("tagline")))
	rec.Set("bio_es", sanitize.HTML(e.Request.PostFormValue("bio_es")))
	rec.Set("bio_en", sanitize.HTML(e.Request.PostFormValue("bio_en")))
	rec.Set("email", strings.TrimSpace(e.Request.PostFormValue("email")))
	rec.Set("phone", strings.TrimSpace(e.Request.PostFormValue("phone")))
	rec.Set("address", strings.TrimSpace(e.Request.PostFormValue("address")))
	rec.Set("show_address", e.Request.PostFormValue("show_address") == "1")
	rec.Set("instagram", strings.TrimSpace(e.Request.PostFormValue("instagram")))
	rec.Set("facebook", strings.TrimSpace(e.Request.PostFormValue("facebook")))
	rec.Set("youtube", strings.TrimSpace(e.Request.PostFormValue("youtube")))
	rec.Set("vimeo", strings.TrimSpace(e.Request.PostFormValue("vimeo")))
	rec.Set("reel_url", strings.TrimSpace(e.Request.PostFormValue("reel_url")))
	if files := uploadedFile(e.Request, "profile_image"); files != nil {
		rec.Set("profile_image", files)
	}
	if files := uploadedFile(e.Request, "hero_image"); files != nil {
		rec.Set("hero_image", files)
	}
	if files := uploadedFiles(e.Request, "hero_images"); len(files) > 0 {
		rec.Set("+hero_images", files)
	}
	if err := app.Save(rec); err != nil {
		f := settingsToForm(rec)
		f.Error = "Error guardando: " + err.Error()
		e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
		return views.AdminSettingsPage(user.GetString("email"), f).Render(context.Background(), e.Response)
	}
	f := settingsToForm(rec)
	f.Success = "Configuración guardada."
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminSettingsPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func settingsToForm(rec *core.Record) views.AdminSettingsForm {
	f := views.AdminSettingsForm{
		SiteName:  rec.GetString("site_name"),
		TagLine:   rec.GetString("tagline"),
		BioES:     rec.GetString("bio_es"),
		BioEN:     rec.GetString("bio_en"),
		Email:     rec.GetString("email"),
		Phone:     rec.GetString("phone"),
		Address:   rec.GetString("address"),
		Instagram: rec.GetString("instagram"),
		Facebook:  rec.GetString("facebook"),
		YouTube:   rec.GetString("youtube"),
		Vimeo:       rec.GetString("vimeo"),
		ReelURL:     rec.GetString("reel_url"),
		ShowAddress: rec.GetBool("show_address"),
	}
	for _, n := range rec.GetStringSlice("hero_images") {
		f.HeroImages = append(f.HeroImages, views.AdminImage{
			URL:      "/api/files/site_settings/" + rec.Id + "/" + n,
			Filename: n,
		})
	}
	if img := rec.GetString("profile_image"); img != "" {
		f.ProfileImage = "/api/files/site_settings/" + rec.Id + "/" + img
	}
	if img := rec.GetString("hero_image"); img != "" {
		f.HeroImage = "/api/files/site_settings/" + rec.Id + "/" + img
	}
	return f
}

// ==================== FORM HELPERS ====================

func readSectionForm(r *http.Request, isNew bool) views.AdminSectionForm {
	order, _ := strconv.Atoi(r.PostFormValue("sort_order"))
	return views.AdminSectionForm{
		IsNew:       isNew,
		Name:        strings.TrimSpace(r.PostFormValue("name")),
		Slug:        strings.TrimSpace(r.PostFormValue("slug")),
		Description: r.PostFormValue("description"),
		SortOrder:   order,
		Active:      r.PostFormValue("active") == "1",
	}
}

func renderSectionForm(e *core.RequestEvent, user *core.Record, f views.AdminSectionForm, errMsg string) error {
	f.Error = errMsg
	e.Response.WriteHeader(http.StatusBadRequest)
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminSectionFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func readWorkForm(r *http.Request, isNew bool) views.AdminWorkForm {
	order, _ := strconv.Atoi(r.PostFormValue("sort_order"))
	return views.AdminWorkForm{
		IsNew:       isNew,
		SectionID:   r.PostFormValue("section"),
		Title:       strings.TrimSpace(r.PostFormValue("title")),
		Slug:        strings.TrimSpace(r.PostFormValue("slug")),
		Year:        strings.TrimSpace(r.PostFormValue("year")),
		Role:        strings.TrimSpace(r.PostFormValue("role")),
		Description: r.PostFormValue("description"),
		Credits:     r.PostFormValue("credits"),
		SortOrder:   order,
		Active:      r.PostFormValue("active") == "1",
		Featured:    r.PostFormValue("featured") == "1",
	}
}

func renderWorkForm(e *core.RequestEvent, user *core.Record, f views.AdminWorkForm, errMsg string) error {
	f.Error = errMsg
	e.Response.WriteHeader(http.StatusBadRequest)
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminWorkFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func readPressForm(r *http.Request, isNew bool) views.AdminPressForm {
	order, _ := strconv.Atoi(r.PostFormValue("sort_order"))
	return views.AdminPressForm{
		IsNew:       isNew,
		Title:       strings.TrimSpace(r.PostFormValue("title")),
		Publication: strings.TrimSpace(r.PostFormValue("publication")),
		URL:         strings.TrimSpace(r.PostFormValue("url")),
		Date:        strings.TrimSpace(r.PostFormValue("date")),
		Excerpt:     r.PostFormValue("excerpt"),
		SortOrder:   order,
		Active:      r.PostFormValue("active") == "1",
	}
}

func renderPressForm(e *core.RequestEvent, user *core.Record, f views.AdminPressForm, errMsg string) error {
	f.Error = errMsg
	e.Response.WriteHeader(http.StatusBadRequest)
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return views.AdminPressFormPage(user.GetString("email"), f).Render(context.Background(), e.Response)
}

func uploadedFile(r *http.Request, field string) *filesystem.File {
	if r.MultipartForm == nil {
		return nil
	}
	headers := r.MultipartForm.File[field]
	if len(headers) == 0 || headers[0].Size == 0 {
		return nil
	}
	f, err := filesystem.NewFileFromMultipart(headers[0])
	if err != nil {
		return nil
	}
	return f
}

func uploadedFiles(r *http.Request, field string) []*filesystem.File {
	if r.MultipartForm == nil {
		return nil
	}
	headers := r.MultipartForm.File[field]
	out := make([]*filesystem.File, 0, len(headers))
	for _, h := range headers {
		if h.Size == 0 {
			continue
		}
		f, err := filesystem.NewFileFromMultipart(h)
		if err != nil {
			continue
		}
		out = append(out, f)
	}
	return out
}
