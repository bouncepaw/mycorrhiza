package views

import (
	"github.com/bouncepaw/mycorrhiza/hyphae/categories"
	"github.com/bouncepaw/mycorrhiza/util"
	"html/template"
	"log"
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
				<label for="_cat-name">
				<input type="text">
				<input type="submit" value="Add to category">
			</form>
		</li>
	</ul>
</aside>`

var categoriesCardT *template.Template

func init() {
	categoriesCardT = template.Must(template.
		New("category card").
		Funcs(template.FuncMap{
			"beautifulName": util.BeautifulName,
		}).
		Parse(categoriesCardTmpl))
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
