package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/gbolo/vsummary/common"
	"github.com/spf13/viper"
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

func handlerUiEsxi(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiEsxi")

	// output the page
	writeSummaryPage(w, &esxiView)

	return
}

func handlerUiPortgroup(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiPortgroup")

	// output the page
	writeSummaryPage(w, &portgroupView)

	return
}

func handlerUiDatastore(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiDatastore")

	// output the page
	writeSummaryPage(w, &datastoreView)

	return
}

// return list of all template filenames
func findAllTemplates() (templateFiles []string, err error) {

	files, err := ioutil.ReadDir(viper.GetString("server.templates_dir"))
	if err != nil {
		return
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".gohtml") {
			templateFiles = append(templateFiles, viper.GetString("server.templates_dir")+"/"+filename)
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
		Funcs(sprig.TxtFuncMap()).
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

// TODO: consolidate this with writeSummaryPage
func writePollerPage(w http.ResponseWriter, t string) {

	// read in all templates
	templateFiles, err := findAllTemplates()

	if err != nil {
		fmt.Fprintf(w, "Error reading template(s). See logs")
		log.Errorf("error geting template files: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// parse and add function to templates
	templates, err := template.New(t).
		Funcs(sprig.TxtFuncMap()).
		ParseFiles(templateFiles...)

	if err == nil {
		pollers, _ := backend.GetPollers()
		execErr := templates.ExecuteTemplate(w, t, UiView{Title: "vCenter Pollers", Pollers: pollers, AjaxEndpoint: common.EndpointPoller})
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

// TODO: consolidate this with writeSummaryPage
func writePollerEditPage(w http.ResponseWriter, t string, id string) {

	// read in all templates
	templateFiles, err := findAllTemplates()

	if err != nil {
		fmt.Fprintf(w, "Error reading template(s). See logs")
		log.Errorf("error geting template files: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// parse and add function to templates
	templates, err := template.New(t).
		Funcs(sprig.TxtFuncMap()).
		ParseFiles(templateFiles...)

	if err == nil {
		poller, errPoller := backend.SelectPoller(id)
		if err != nil {
			fmt.Fprintf(w, "Error finding poller. See logs")
			log.Errorf("template execute error: %s", errPoller)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		pollers := []common.Poller{poller}
		execErr := templates.ExecuteTemplate(w, t, UiView{Title: "vCenter Pollers", Pollers: pollers, AjaxEndpoint: common.EndpointPoller})
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
