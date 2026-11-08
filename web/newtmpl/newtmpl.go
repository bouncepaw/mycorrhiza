package newtmpl

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
	"github.com/bouncepaw/mycorrhiza/web/views"
)

var base = template.Must(template.ParseFS(views.FS, "base.html"))

type Page struct {
	Template *template.Template
}

func NewPage(tmpls ...string) *Page {
	must := template.Must
	tmpl := must(must(
		base.Clone()).
		Funcs(template.FuncMap{
			"beautifulName":     util.BeautifulName,
			"wikiName":          func() string { return cfg.WikiName },
			"userHypha":         func() string { return cfg.UserHypha },
			"useAuth":           func() bool { return cfg.UseAuth },
			"allowRegistration": func() bool { return cfg.AllowRegistration },
			"inc":               func(i int) int { return i + 1 },
			"base": func(hyphaName string) string {
				parts := strings.Split(hyphaName, "/")
				return parts[len(parts)-1]
			},
			"beautifulLink": func(hyphaName string) template.HTML {
				return template.HTML(
					fmt.Sprintf(
						`<a href="/hypha/%s">%s</a>`, hyphaName, hyphaName))
			},
		}).
		ParseFS(views.FS, append([]string{"shared.html"}, tmpls...)...))

	return &Page{Template: tmpl}
}

func (p *Page) RenderTo(meta viewdata.Meta, data map[string]any) error {
	var common BaseData
	common.setBaseValuesFrom(meta)
	data["Meta"] = common.Meta
	data["HeadElements"] = common.HeadElements
	data["BodyAttributes"] = common.BodyAttributes
	data["CommonScripts"] = common.CommonScripts
	data["EditScripts"] = common.EditScripts
	data["HeaderLinks"] = common.HeaderLinks
	data["UseAuth"] = common.UseAuth

	return p.Template.ExecuteTemplate(meta.W, "page", data)
}

type BaseData struct {
	Meta           viewdata.Meta
	HeadElements   []template.HTML
	BodyAttributes map[string]string
	CommonScripts  []string
	EditScripts    []string
	HeaderLinks    []viewdata.HeaderLink
	UseAuth        bool
	Addr           string
}

func (bd *BaseData) setBaseValuesFrom(meta viewdata.Meta) {
	bd.Meta = meta
	bd.HeadElements = meta.HeadElements
	bd.BodyAttributes = meta.BodyAttributes
	bd.CommonScripts = cfg.CommonScripts
	bd.EditScripts = cfg.EditScripts
	bd.HeaderLinks = viewdata.HeaderLinks
	bd.UseAuth = cfg.UseAuth
}

func (p *Page) RenderDataTo(meta viewdata.Meta, data interface {
	setBaseValuesFrom(viewdata.Meta)
}) error {
	data.setBaseValuesFrom(meta)
	return p.Template.ExecuteTemplate(meta.W, "page", data)
}
