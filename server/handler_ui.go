package server

import (
	"text/template"
	"net/http"
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

	templates, err := template.ParseFiles("www/templates/header.gohtml", "www/templates/navigation.gohtml", "www/templates/index.gohtml")
	if err == nil {
		execErr := templates.ExecuteTemplate(w, "index", ui)
		if execErr != nil {
			log.Errorf("exec template error: %s", execErr)
		}

	} else {
		log.Errorf("template error: %s", err)
	}


	return
}
