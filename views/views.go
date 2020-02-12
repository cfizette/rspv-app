package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	LayoutDir         = "views/layouts/"
	TemplateExtension = ".html"
	TemplateDir       = "views/"
)

type View struct {
	Template *template.Template
	Layout   string
}

// NewView creates View struct from the given layout and files.
func NewView(layout string, files ...string) *View {
	addTemplateExt(files)
	addTemplatePath(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}

}

// Render renders the view to the ResponseWriter
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// ServeHTTP implements the Handler interface
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings representing file paths
// for templates and prepends the TemplateDir directory to each string
// in the slice
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings representing file paths
// for templates and appends the TemplateExtension to each string in
// the slice
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExtension
	}
}
