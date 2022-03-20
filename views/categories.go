package views

import (
	"github.com/bouncepaw/mycorrhiza/hyphae/categories"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

const categoriesCardTmpl = `{{$hyphaName := .HyphaName
}}<aside class="layout-card categories-card">
	<h2 class="layout-card__title">Categories</h2>
	<ul class="categories-card__entries">
	{{range .Categories}}
		<li class="categories-card__entry">
			<a class="categories-card__link" href="/category/{{.}}">{{beautifulName .}}</a>
			<form method="POST" action="/remove-from-category" class="categories-card__remove-form">
				<input type="hidden" name="cat" value="{{.}}">
				<input type="hidden" name="hypha" value="{{$hyphaName}}">
				<input type="submit" value="X">
			</form>
		</li>
	{{end}}
		<li class="categories-card__entry categories-card__add-to-cat">
			<form method="POST" action="/add-to-category" class="categories-card__add-form">
				<input type="text" name="cat" id="_cat-name">
				<input type="hidden" name="hypha" value="{{$hyphaName}}">
				<input type="submit" value="Add to category">
			</form>
		</li>
	</ul>
</aside>`

const categoryPageTmpl = `{{$catName := .CatName
}}<main class="main-width category">
	<h1>Category <i>{{$catName}}</i></h1>
{{if len .Hyphae}}
	<p>This page lists all hyphae in the category.</p>
{{else}}
	<p>This category has no hyphae.</p>
{{end}}
	<ul class="category__entries">
	{{range .Hyphae}}
		<li class="category__entry">
			<a class="wikilink" href="/hypha/{{.}}">{{beautifulName .}}</a>
		</li>
	{{end}}
		<li class="category__entry category__add-to-cat">
			<form method="POST" action="/add-to-category" class="category__add-form">
				<input type="text" name="hypha" id="_hypha-name" placeholder="Hypha name">
				<input type="hidden" name="cat" value="{{$catName}}">
				<input type="submit" value="Add hypha to category">
			</form>
		</li>
	</ul>
</main>`

var (
	categoriesCardT *template.Template
	categoryPageT   *template.Template
)

func init() {
	categoriesCardT = template.Must(template.
		New("category card").
		Funcs(template.FuncMap{
			"beautifulName": util.BeautifulName,
		}).
		Parse(categoriesCardTmpl))
	categoryPageT = template.Must(template.
		New("category page").
		Funcs(template.FuncMap{
			"beautifulName": util.BeautifulName,
		}).
		Parse(categoryPageTmpl))
}

func categoryCardHTML(hyphaName string) string {
	var buf strings.Builder
	err := categoriesCardT.Execute(&buf, struct {
		HyphaName  string
		Categories []string
	}{
		hyphaName,
		categories.WithHypha(hyphaName),
	})
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

func CategoryPageHTML(w io.Writer, rq *http.Request, catName string) {
	var buf strings.Builder
	err := categoryPageT.Execute(&buf, struct {
		CatName string
		Hyphae  []string
	}{
		catName,
		categories.Contents(catName),
	})
	if err != nil {
		log.Println(err)
	}
	io.WriteString(w, BaseHTML(
		"Category "+util.BeautifulName(catName),
		buf.String(),
		l18n.FromRequest(rq),
		user.FromRequest(rq),
	))
}
