package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// 自定义模板函数 action
func humanDate(t time.Time) string {
	// 格式化time
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// 解析模板并缓存模板对象
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		// 提取*.html部分
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		// ts中是一个模板集，存有所有模板名字的slice，通过
		// ExecuteTemplate(w, "template—name", date)
		// 可以渲染对应名字的模板
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

// Test
// func main() {
// fmt.Printf("%v\n", page)
//
// fmt.Printf("%v\n", name)
// 	newTemplateCache()
// }
