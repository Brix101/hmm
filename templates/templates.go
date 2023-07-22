package templates

import (
	"embed"
	"html/template"
)

//go:embed views/*.html
var tmplFS embed.FS

//go:embed static
var staticFs embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	funcMap := template.FuncMap{
		"inc": inc,
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(tmplFS, "views/*.html"))
	return &Template{
		templates: templates,
	}
}
