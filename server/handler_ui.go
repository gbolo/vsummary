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

// default handler for landing page
func handlerUiIndex(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handleUiIndex")

	// for now just return a 302 to the VM endpoint
	http.Redirect(w, req, "/ui/virtualmachines", http.StatusTemporaryRedirect)
}

func handlerUiVirtualmachines(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiVirtualmachines")

	// output the page
	writeSummaryPage(w, &virtualMachineView)
}

func handlerUiDatacenters(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handleUiDatacenters")

	// page details
	uiview := UiView{
		Title:        "Datacenters",
		AjaxEndpoint: "/api/dt/datacenters",
		Table: map[string]string{
			"name":       "Name",
			"vcenter_id": "vCenter Id",
		},
	}

	// output the page
	writeSummaryPage(w, &uiview)

	return
}

func handlerUiEsxi(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiEsxi")

	// page details
	uiview := UiView{
		Title:        "ESXi",
		AjaxEndpoint: "/api/dt/esxi",
		Table: map[string]string{
			"name": "Name",
		},
	}

	// output the page
	writeSummaryPage(w, &uiview)

	return
}

// return list of all template filenames
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

// write generic vsummary table page
func writeSummaryPage(w http.ResponseWriter, uiview *UiView) {

	// read in all templates
	templateFiles, err := findAllTemplates()

	if err != nil {
		fmt.Fprintf(w, "Error reading template(s). See logs")
		log.Errorf("error geting template files: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// parse and add function to templates
	templates, err := template.New("index").
		Funcs(template.FuncMap{"StringsJoin": strings.Join}).
		ParseFiles(templateFiles...)

	if err == nil {
		execErr := templates.ExecuteTemplate(w, "index", uiview)
		if execErr != nil {
			fmt.Fprintf(w, "Error executing template(s). See logs")
			log.Errorf("template execute error: %s", execErr)
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


// sets the values of datatables json columns definition
func setDtColumns(uiview *UiView) {

	// prepare slices for templates
	for colName, trName := range uiview.Table {
		uiview.DtColumns = append(uiview.DtColumns, fmt.Sprintf(
			"{ \"data\": \"%s\", \"name\": \"%s\", \"title\": \"%s\" }",
			colName,
			colName,
			trName,
		))
	}
}