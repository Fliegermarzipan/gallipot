package template

import (
	//"net/http"
	"github.com/Fliegermarzipan/gallipot/data"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type FrontendData struct {
	//Request *http.Request
	LoggedIn bool
	Page     string
	User     *data.User
}

var (
	indexTemplate *template.Template
)

func LoadTemplates() {
	funcMap := template.FuncMap{
		"split": strings.Split,
		"contains": func(s []string, e string) bool {
			for _, a := range s {
				if a == e {
					return true
				}
			}
			return false
		},
		"combine": func(s ...[]string) []string {
			ret := []string{}
			for i := range s {
				ret = append(ret, s[i]...)
			}
			return ret
		},
		"list": func(s ...interface{}) []interface{} {
			return s
		},
		"lower": strings.ToLower,
		"title": strings.Title,
		"exists": func(s string) bool {
			if _, err := os.Stat(s); err == nil {
				return true
			} else if os.IsNotExist(err) {
				return false
			}
			log.Print("file exists fucked up")
			return false
		},
	}
	funcMap["include"] = func(s string, d interface{}) interface{} {
		ext := filepath.Ext(s)

		var buf strings.Builder
		err := template.Must(template.New(path.Base(s)).Funcs(funcMap).ParseFiles(s)).Execute(&buf, d)
		if err != nil {
			log.Print(err)
		}
		rendered := buf.String()

		if ext == ".html" {
			return template.HTML(rendered)
		} else if ext == ".css" {
			return template.CSS(rendered)
		}
		return rendered
	}

	_, caller, _, _ := runtime.Caller(0)
	caller, err := filepath.EvalSymlinks(caller)
	if err != nil {
		log.Fatal(err)
	}
	templateDir := filepath.Dir(caller)
	indexTemplate = template.Must(template.New("").Funcs(funcMap).
		ParseGlob(filepath.Join(templateDir, "*.html"))).Lookup("index.html")
}

func GetTemplates() *template.Template {
	if indexTemplate == nil {
		LoadTemplates()
	}
	return indexTemplate
}
