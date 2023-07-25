package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed views/layouts
//go:embed views
var tmplFS embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	funcMap := template.FuncMap{
		// "inc": inc,
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(tmplFS, "views/*.html", "views/layouts/*.html"))
	return &Template{
		templates: templates,
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	fmt.Println("++++++++++++++++", "views/"+name)
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(tmplFS, "views/"+name))
	return tmpl.ExecuteTemplate(w, name, data)
}
