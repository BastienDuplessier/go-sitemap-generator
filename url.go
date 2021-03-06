package sitemap

const (
	ChangeFreqAlways  = "always"
	ChangeFreqHourly  = "hourly"
	ChangeFreqDaily   = "daily"
	ChangeFreqWeekly  = "weekly"
	ChangeFreqMonthly = "monthly"
	ChangeFreqYearly  = "yearly"
	ChangeFreqNever   = "never"
)

type SitemapURL interface {
	SitemapLoc() string
	SitemapChangeFreq() string
	SitemapLastMod() string
	SitemapPriority() string
	SitemapImages() []Image
	SitemapAlternateLinks() []AlternateLink
}

type URL struct {
	Loc            string
	ChangeFreq     string
	LastMod        string
	Priority       string
	Images         []Image
	AlternateLinks []AlternateLink
}

func (u URL) SitemapLoc() string {
	return u.Loc
}

func (u URL) SitemapChangeFreq() string {
	return u.ChangeFreq
}

func (u URL) SitemapLastMod() string {
	return u.LastMod
}

func (u URL) SitemapPriority() string {
	return u.Priority
}

func (u URL) SitemapImages() []Image {
	return u.Images
}

func (u URL) SitemapAlternateLinks() []AlternateLink {
	return u.AlternateLinks
}

type AlternateLink struct {
	Url  string
	Lang string
}
