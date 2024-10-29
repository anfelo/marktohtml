package main

import (
	"embed"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed html/*.html
var files embed.FS

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}

func main() {
	templates := make(map[string]*template.Template)
	templates["home"] = template.Must(
		template.New("html/layout.html").ParseFS(files, "html/layout.html", "html/home.html"),
	)

	e := echo.New()
	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:CSRF",
	}))

	e.Static("/static", "public")

	e.GET("/", Home)
	e.POST("/markdown-to-html", MarkdownToHTML)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func Home(c echo.Context) error {
	type data struct {
		CSRF string
	}

	return c.Render(http.StatusOK, "home", data{
		CSRF: c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
	})
}

func MarkdownToHTML(c echo.Context) error {
	mdStr := c.FormValue("markdown")

    mdBytes := []byte(mdStr)
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdBytes)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

    result := string(markdown.Render(doc, renderer))

	return c.HTML(http.StatusOK, result)
}
