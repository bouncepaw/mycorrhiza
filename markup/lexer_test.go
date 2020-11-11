package markup

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

// TODO: move test markup docs to files, perhaps? These strings sure are ugly
func TestLex(t *testing.T) {
	check := func(name, content string, expectedAst []Line) {
		if ast := lex(name, content); !reflect.DeepEqual(ast, expectedAst) {
			if len(ast) != len(expectedAst) {
				t.Error("Expected and generated AST length of", name, "do not match. Printed generated AST.")
				for _, l := range ast {
					fmt.Printf("%d: %s\n", l.id, l.contents)
				}
				return
			}
			for i, e := range ast {
				if !reflect.DeepEqual(e, expectedAst[i]) {
					t.Error(fmt.Sprintf("Expected: %q\nGot:%q", expectedAst[i], e))
				}
			}
		}
	}
	contentsB, err := ioutil.ReadFile("testdata/test.myco")
	if err != nil {
		t.Error("Could not read test markup file!")
	}
	contents := string(contentsB)
	check("Apple", contents, []Line{
		{1, "<h1 id='1'>1</h1>"},
		{2, "<h2 id='2'>2</h2>"},
		{3, "<h3 id='3'>3</h3>"},
		{4, "<blockquote id='4'>quote</blockquote>"},
		{5, `<ul id='5'>
	<li>li 1</li>
	<li>li 2</li>
</ul>`},
		{6, "<p id='6'>text</p>"},
		{7, "<p id='7'>more text</p>"},
		{8, `<p><a id='8' class='rocketlink wikilink_internal' href="/page/pear">some link</a></p>`},
		{9, `<ul id='9'>
	<li>lin&#34;+</li>
</ul>`},
		{10, `<pre id='10' alt='alt text goes here' class='codeblock'><code>=&gt; preformatted text
where markup is not lexed</code></pre>`},
		{11, `<p><a id='11' class='rocketlink wikilink_internal' href="/page/linking">linking</a></p>`},
		{12, "<p id='12'>text</p>"},
		{13, `<pre id='13' alt='' class='codeblock'><code>()
/\</code></pre>`},
		{14, Transclusion{"apple", 1, 3}},
		{15, Img{
			hyphaName: "Apple",
			inDesc:    false,
			entries: []imgEntry{
				{"/binary/hypha1", "", "", ""},
				{"/binary/hypha2", "", "", ""},
				{"/binary/hypha3", "60", "", ""},
				{"/binary/hypha4", "", "", " line1\nline2\n"},
				{"/binary/hypha5", "", "", "\nstate of minnesota\n"},
			},
		}},
	})
}
