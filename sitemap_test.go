package sitemap

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

const tmpDir = `testdata`

func init() {
	t, err := filepath.Abs(tmpDir)
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(t)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(t, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func TestGenerator(t *testing.T) {
	o := Options{
		MaxURLs:  2,
		Dir:      tmpDir,
		Filename: "a",
		BaseURL:  "http://example.com/",
	}

	g := New(o)
	require.NoError(t, g.Open())
	require.NoError(t, g.Add(URL{Loc: "test1"}))
	require.NoError(t, g.Add(URL{Loc: "test2"}))
	require.NoError(t, g.Close())

	o.Filename = "b"
	g = New(o)
	require.NoError(t, g.Open())
	require.NoError(t, g.Add(URL{Loc: "test1"}))
	require.NoError(t, g.Add(URL{Loc: "test2"}))
	require.NoError(t, g.Add(URL{Loc: "test3"}))
	require.NoError(t, g.Add(URL{Loc: "test4"}))
	require.NoError(t, g.Add(URL{Loc: "test5"}))
	require.NoError(t, g.Close())

	files, _ := ioutil.ReadDir(g.opt.Dir)
	require.Equal(t, 5, len(files))
	require.Equal(t, "a.xml", files[0].Name())
	require.Equal(t, "b-1.xml", files[1].Name())
	require.Equal(t, "b-2.xml", files[2].Name())
	require.Equal(t, "b-3.xml", files[3].Name())
	require.Equal(t, "b.xml", files[4].Name())
}

func TestParamChecks(t *testing.T) {
	var (
		g   *Generator
		err error
	)
	opt := Options{
		MaxFileSize: 2,
		MaxURLs:     -1,
		BaseURL:     "/",
		Dir:         tmpDir,
	}
	g = New(opt)
	err = g.Open()
	require.Error(t, err)

	opt.MaxFileSize = len(header+baseOpen) + len(baseClose) + 10
	g = New(opt)
	err = g.Open()
	require.Error(t, err)

	opt.MaxURLs = 2
	g = New(opt)
	err = g.Open()
	require.NoError(t, err)
}

func TestHeaderChecks(t *testing.T) {
	var g *Generator
	checks := [3]string{
		`xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`,
		`xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"`,
		`xmlns:xhtml="http://www.w3.org/1999/xhtml"`,
	}

	ns := map[string]string{
		"xmlns:xhtml": "http://www.w3.org/1999/xhtml",
		"xmlns:image": "http://www.google.com/schemas/sitemap-image/1.1",
	}
	opt := Options{
		Dir:      tmpDir,
		Filename: "c",
		BaseURL:  "http://example.com/",
		XMLns:    ns,
	}
	g = New(opt)
	require.NoError(t, g.Open())
	require.NoError(t, g.Add(URL{Loc: "test1"}))
	require.NoError(t, g.Add(URL{Loc: "test2"}))
	require.NoError(t, g.Close())

	file, _ := ioutil.ReadFile(tmpDir + "/c.xml")
	for _, check := range checks {
		regString := "<urlset.*" + check + ".*>.*</urlset>"
		t.Log("Testing presence of: " + check)
		res, err := regexp.Match(regString, file)
		require.NoError(t, err)
		require.True(t, res)
	}
}

func TestInternals(t *testing.T) {
	var (
		g   *Generator
		err error
	)
	opt := Options{
		MaxFileSize: len(header+baseOpen) - 2 + len(baseClose) + 10,
		MaxURLs:     2,
		BaseURL:     "/",
		Dir:         tmpDir,
	}
	g = New(opt)
	err = g.Open()
	require.NoError(t, err)
	require.True(t, g.canFit(10))
	require.False(t, g.canFit(11))

	n1 := g.formatURLNode(URL{Loc: "test1"})
	require.Equal(t, `<url><loc>test1</loc></url>`, n1)

	n2 := g.formatURLNode(URL{Loc: "test2", ChangeFreq: ChangeFreqDaily})
	require.Equal(t, `<url><loc>test2</loc><changefreq>daily</changefreq></url>`, n2)
}

func TestInternalsImages(t *testing.T) {
	var (
		g   *Generator
		err error
	)
	ns := map[string]string{
		"xmlns:image": "http://www.google.com/schemas/sitemap-image/1.1",
	}
	opt := Options{
		Dir:      tmpDir,
		Filename: "c",
		BaseURL:  "http://example.com/",
		XMLns:    ns,
	}

	g = New(opt)
	err = g.Open()
	require.NoError(t, err)

	img1 := Image{Loc: "loc1", Title: "img1"}
	img2 := Image{
		Loc:     "loc2",
		Caption: "image2",
		GeoLoc:  "somewhere",
		Title:   "img2",
		License: "url",
	}

	n1 := g.formatURLNode(URL{Loc: "test1", Images: []Image{img1}})
	require.Equal(t, "<url>"+
		"<loc>test1</loc>"+
		"<image:image>"+
		"<image:loc>loc1</image:loc>"+
		"<image:title>img1</image:title>"+
		"</image:image>"+
		"</url>", n1)

	n2 := g.formatURLNode(URL{Loc: "test2", Images: []Image{img1, img2}})
	require.Equal(t, "<url>"+
		"<loc>test2</loc>"+
		"<image:image>"+
		"<image:loc>loc1</image:loc>"+
		"<image:title>img1</image:title>"+
		"</image:image>"+
		"<image:image>"+
		"<image:loc>loc2</image:loc>"+
		"<image:caption>image2</image:caption>"+
		"<image:geo_location>somewhere</image:geo_location>"+
		"<image:title>img2</image:title>"+
		"<image:license>url</image:license>"+
		"</image:image>"+
		"</url>", n2)
}
