package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gomarkdown/markdown"
)

type Page struct {
	Title   string
	Nav     template.HTML
	Content template.HTML
	URL     string
	Posts   []Page
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func generate_post_pages() []Page {
	var files []string
	var posts []Page
	root := "./posts"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != "./posts" {
			files = append(files, path)
		}
		return nil
	})

	nav, err := ioutil.ReadFile("components/nav.html")
	check(err)

	t, err := template.ParseFiles("templates/post-body.html")
	check(err)

	for _, path := range files {
		content, err := ioutil.ReadFile(path)
		check(err)

		md := []byte(content)
		md_output := markdown.ToHTML(md, nil, nil)

		m := regexp.MustCompile(`posts\/([a-z\-]*).md`)
		slug := m.FindStringSubmatch(path)[1]

		postPage := Page{
			Title:   slug,
			Nav:     template.HTML(string(nav)),
			Content: template.HTML(string(md_output)),
			URL:     "posts/" + slug + ".html",
		}

		posts = append(posts, postPage)

		fmt.Println("Generating " + slug + ".html ......")
		var finalOutput bytes.Buffer

		err = t.Execute(&finalOutput, postPage)
		check(err)

		ferr := ioutil.WriteFile("./dist/" + postPage.URL, []byte(finalOutput.String()), 0644)
		check(ferr)
	}
	return posts
}

func generate_front_page(posts []Page) {
	var finalOutput bytes.Buffer

	t, err := template.ParseFiles("pages/index.html")
	check(err)

	nav, err := ioutil.ReadFile("components/nav.html")
	check(err)

	frontPage := Page{
		Title: "Dunedain Dev",
		Nav:   template.HTML(string(nav)),
		URL:   "//index.html",
		Posts: posts,
	}

	err = t.Execute(&finalOutput, frontPage)
	check(err)

	ferr := ioutil.WriteFile("./dist" + frontPage.URL, []byte(finalOutput.String()), 0644)
	check(ferr)
}

func main() {
	posts := generate_post_pages()
	generate_front_page(posts)
}
