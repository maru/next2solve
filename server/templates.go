// Next problem to solve
// https://github.com/maru/next2solve
//
// HTTP handler

package server

import (
	"fmt"
	"html/template"
	"net/http"
	"next2solve/problems"
)

type TemplateData struct {
	UsernameError   string
	UserID          string
	Username        string
	Problems        []problems.ProblemInfo
	IsOrderStar     bool
	IsOrderCategory bool
	IsOrderLevel    bool
}

var (
	funcMap = template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}
	paths = []string{
		"templates/header.html",
		"templates/footer.html",
		"templates/form.html",
		"templates/index.html",
		"templates/lucky.html",
		"templates/problems.html",
	}
	tmpl *template.Template
)

// Render page using a template with data
func renderPage(w http.ResponseWriter, name string, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		fmt.Fprintf(w, "Error %v", err)
	}
}

// Set template data
func setTemplateUserData(username string) (TemplateData, error) {
	// Check if username is valid
	userid, err := problems.GetUserID(username)
	if err != nil {
		return TemplateData{UsernameError: err.Error(), Username: username}, err
	}
	// Set user data
	return TemplateData{UserID: userid, Username: username}, nil
}

func LoadTemplates(dir ...string) {
	// This function may be called from tests...
	if dir != nil {
		for i, file := range paths {
			paths[i] = dir[0] + "/" + file
		}
	}
	tmpl = template.Must(template.New("").Funcs(funcMap).ParseFiles(paths...))
}
