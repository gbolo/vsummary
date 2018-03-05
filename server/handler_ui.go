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
	Title        string
	AjaxEndpoint string
	TableHeaders []string
}

// default handler for landing page
func handlerUiIndex(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handleUiIndex")

	// for now just return a 302 to the VM endpoint
	http.Redirect(w, req, "/virtualmachines", http.StatusTemporaryRedirect)
}

func handlerUiVirtualmachines(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handlerUiVirtualmachines")

	uiview := UiView{
		"Virtualmachines",
		"/api/dt/virtualmachines",
		[]string{
			"Folder",
			"vCPU",
			"Memory",
			"PowerState",
			"Real GuestOS",
			"Config GuestOS",
			"Version",
			"ConfigChange",
			"ToolsVersion",
			"ToolRunning",
			"Hostname",
			"IP",
			"Cluster",
			"Pool",
			"Datacenter",
			"CpuUsed",
			"HostMemUsed",
			"GuestMemUsed",
			"Uptime",
			"ESXi",
			"ESXiEVC",
			"ESXiStatus",
			"ESXiCPU",
			"vDisks",
			"vNICs",
			"VMX",
			"vCenter",
			"VC-ENV",
		},
	}

	// output the page
	writeSummaryPage(w, &uiview)

	return
}

func handlerUiDatacenters(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handleUiDatacenters")

	// page details
	uiview := UiView{
		"Datacenters",
		"/api/dt/datacenters",
		[]string{
			"Name",
			"EsxiFolderId",
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

	// parse then execute/render templates
	templates, err := template.ParseFiles(templateFiles...)
	if err == nil {
		execErr := templates.ExecuteTemplate(w, "index", uiview)
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
