package views

// ---------- Admin view types ----------

type AdminSectionRow struct {
	ID         string
	Name       string
	Slug       string
	SortOrder  int
	Active     bool
	WorksCount int
}

type AdminSectionForm struct {
	IsNew       bool
	ID          string
	Name        string
	Slug        string
	Description string
	SortOrder   int
	Active      bool
	CoverURL    string
	Error       string
}

type AdminWorkRow struct {
	ID          string
	Title       string
	Slug        string
	Year        string
	Role        string
	SectionName string
	SectionID   string
	ThumbURL    string
	Active      bool
	Featured    bool
	LinksCount  int
}

type AdminWorkForm struct {
	IsNew       bool
	ID          string
	SectionID   string
	SectionName string
	Title       string
	Slug        string
	Year        string
	Role        string
	Description string
	Credits     string
	SortOrder   int
	Active      bool
	Featured    bool
	Images      []AdminImage
	Links       []AdminWorkLink
	Error       string
	Sections    []SectionOption
}

type AdminImage struct {
	URL      string
	Filename string
}

type AdminWorkLink struct {
	ID    string
	Label string
	URL   string
	Kind  string
}

type SectionOption struct {
	ID   string
	Name string
}

type AdminPressRow struct {
	ID          string
	Title       string
	Publication string
	URL         string
	Date        string
	Active      bool
}

type AdminPressForm struct {
	IsNew       bool
	ID          string
	Title       string
	Publication string
	URL         string
	Date        string
	Excerpt     string
	SortOrder   int
	Active      bool
	Error       string
}

type AdminSettingsForm struct {
	SiteName     string
	TagLine      string
	BioES        string
	BioEN        string
	Email        string
	Phone        string
	Address      string
	Instagram    string
	Facebook     string
	YouTube      string
	Vimeo        string
	Website      string
	ReelURL      string
	ProfileImage string
	HeroImage    string
	HeroImages   []AdminImage
	ShowAddress  bool
	Error        string
	Success      string
}
