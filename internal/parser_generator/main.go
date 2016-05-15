package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const goHeader = `package main

func acceptString(input, expected string) (string, bool, error) {
	return input, false, nil
}

func acceptOne(input string, options []string) (string, bool, error) {
	return input, false, nil
}

`

func variants(p string, in []string) [][]string {
	a := [][]string{{p}}

	for i := 0; i < len(in); i++ {
		a = append(a, []string{p, in[i]})

		for n, m := i+1, len(in); n < m; n++ {
			a = append(a, append([]string{p, in[i]}, in[n:m]...))
		}
	}

	return a
}

var (
	htmlFlag   = flag.String("html", "./index.html", "Location of the ECMAScript HTML specification")
	formatFlag = flag.String("format", "ebnf", "Format of output (ebnf, go)")
)

func main() {
	flag.Parse()

	var d *goquery.Document

	if strings.HasPrefix(*htmlFlag, "http") {
		r, err := http.Get(*htmlFlag)
		if err != nil {
			panic(err)
		}

		_d, err := goquery.NewDocumentFromResponse(r)
		if err != nil {
			panic(err)
		}

		d = _d
	} else {
		fd, err := os.Open(*htmlFlag)
		if err != nil {
			panic(err)
		}

		_d, err := goquery.NewDocumentFromReader(fd)
		if err != nil {
			panic(err)
		}

		d = _d
	}

	productions := make(map[string]*production)
	var plist []*production

	d.Find("emu-production[id]").Each(func(i int, ptag *goquery.Selection) {
		id := ptag.AttrOr("id", "")
		if strings.HasPrefix(id, "prod-grammar-notation-") || strings.HasPrefix(id, "prod-asi-rules-") || strings.HasPrefix(id, "prod-annexB-") || ptag.HasClass("inline") {
			return
		}

		p := production{
			name: ptag.AttrOr("name", ""),
		}

		if s, ok := ptag.Attr("params"); ok && s != "" {
			p.params = strings.Split(s, ", ")
		}

		if ptag.Find("emu-oneof").Length() > 0 {
			var a []string
			ptag.Find("emu-t").Each(func(i int, ttag *goquery.Selection) {
				a = append(a, ttag.Text())
			})

			p.rules = append(p.rules, rule{
				tokens: []token{
					token{
						oneof:  true,
						values: a,
					},
				},
			})
		} else {
			ptag.Find("emu-rhs").Each(func(i int, rtag *goquery.Selection) {
				var r rule

				if s, ok := rtag.Attr("constraints"); ok && s != "" {
					for _, c := range strings.Split(s, ", ") {
						r.constraints = append(r.constraints, constraint{
							inverse: c[0] == '~',
							value:   c[1:],
						})
					}
				}

				rtag.ChildrenFiltered("emu-t, emu-nt, emu-gprose").Each(func(i int, etag *goquery.Selection) {
					var t token

					switch etag.Nodes[0].Data {
					case "emu-t":
						t.terminal = true
						t.value = etag.Text()
					case "emu-nt":
						t.value = etag.Find("a").Text()
					case "emu-gprose":
						t.terminal = true
						t.value = "<" + etag.Text() + ">"
					}

					if _, ok := etag.Attr("optional"); ok {
						t.optional = true
					}

					if s, ok := etag.Attr("params"); ok && s != "" {
						t.params = strings.Split(s, ", ")
					}

					r.tokens = append(r.tokens, t)
				})

				p.rules = append(p.rules, r)
			})
		}

		productions[p.name] = &p
		plist = append(plist, &p)
	})

	if *formatFlag == "go" {
		fmt.Printf(goHeader)
	}

	for _, p := range plist {
		if *formatFlag == "go" {
			fmt.Printf("%s\n", p.Go(productions))
		}

		if *formatFlag == "ebnf" {
			for _, s := range p.EBNF() {
				fmt.Printf("%s\n\n", s)
			}
		}
	}
}
