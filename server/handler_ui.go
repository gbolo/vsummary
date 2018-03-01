package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gbolo/vsummary/common"
)

type UiView struct {
	Title string
}

func handlerUiView(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handleUiView")

	ui := UiView{"IndexPage"}

	// read in all templates
	templateFiles, err := findAllTemplates()

	if err != nil {
		fmt.Fprintf(w, "Error reading template(s). See logs")
		log.Errorf("error geting template files: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// parse then execute/render templates
	templates, err := template.ParseFiles(templateFiles...)
	if err == nil {
		execErr := templates.ExecuteTemplate(w, "index", ui)
		if execErr != nil {
			fmt.Fprintf(w, "Error executing template(s). See logs")
			log.Errorf("template execute error: %s", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

	} else {
		fmt.Fprintf(w, "Error parsing template(s). See logs")
		log.Errorf("template parse error: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	return
}

func findAllTemplates() (templateFiles []string, err error) {

	files, err := ioutil.ReadDir("./www/templates")
	if err != nil {
		return
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".gohtml") {
			templateFiles = append(templateFiles, "./www/templates/"+filename)
		}
	}

	return
}
