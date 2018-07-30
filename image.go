package sitemap

type SitemapImage interface {
	SitemapLoc() string
	SitemapCaption() string
	SitemapGeoLoc() string
	SitemapTitle() string
	SitemapLicense() string
}

type Image struct {
	Loc     string
	Caption string
	GeoLoc  string
	Title   string
	License string
}

func (i Image) SitemapLoc() string {
	return i.Loc
}

func (i Image) SitemapCaption() string {
	return i.Caption
}

func (i Image) SitemapGeoLoc() string {
	return i.GeoLoc
}

func (i Image) SitemapTitle() string {
	return i.Title
}

func (i Image) SitemapLicense() string {
	return i.License
}
