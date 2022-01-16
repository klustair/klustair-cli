package klustair

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/klustair/klustair-go/pkg/api"
	"github.com/klustair/klustair-go/pkg/kubeaudit"
	log "github.com/sirupsen/logrus"
)

type Report struct {
	Uid               string `json:"uid"`
	Label             string `json:"title"`
	namespaces        *NamespaceList
	objectsList       *ObjectsList
	kubeauditReports  [][]kubeaudit.KubeauditReport
	targetslist       Targetslist
	reportSummary     *ReportSummary
	containerHasImage []*ContainerHasImage
}

func (r *Report) Init(label string, whitelist []string, blacklist []string, trivy bool, kubeauditAuditors []string) {
	r.Uid = uuid.New().String()
	r.Label = label

	ns := new(NamespaceList)
	ns.Init(whitelist, blacklist)
	r.namespaces = ns

	o := new(ObjectsList)
	o.Init(r.Uid, r.namespaces)
	r.objectsList = o

	// recate a report summary
	rs := new(ReportSummary)
	rs.Init()
	rs.NamespacesTotal = ns.Total
	rs.NamespacesChecked = ns.Checked
	rs.Pods = len(o.pods)
	rs.Containers = len(o.containers)
	rs.Images = len(o.uniqueImages)

	// run kubeaudit scans if enabled
	if len(kubeauditAuditors) > 0 && kubeauditAuditors[0] != "" {
		k := new(kubeaudit.Auditor)
		k.Klustair.ReportUid = r.Uid
		k.SetConfig(kubeauditAuditors)

		for _, namespace := range r.namespaces.Namespaces {
			log.Infof("Kubeaudit on namespace: %+v", namespace.Name)
			report := k.Audit(namespace.Name)
			k.Klustair.NamespaceUid = namespace.Uid

			r.kubeauditReports = append(r.kubeauditReports, report)
		}
	}

	// run trivy scans if enabled
	if trivy {
		r.targetslist = o.ScanImages()
		//os.Exit(0)
		rs.sumVulnSummary(r.objectsList.uniqueImages)
	}

	r.containerHasImage = o.linkImagesToContainers(trivy)

	r.reportSummary = rs

	log.Debugf("REPORT SUMMARY: %+v\n", r.reportSummary)
}

func NewReport(opt Options) *Report {
	r := new(Report)
	r.Init(opt.Label, opt.Namespaces, opt.NamespacesBlacklist, opt.Trivy, opt.KubeAudit)
	return r
}

func (r *Report) Send(opt Options) error {
	fmt.Println("SEND REPORT ------------------>>")
	apiClient := api.NewApiClient(opt.Apihost, opt.Apitoken)

	////////////////////////////////////////////////////////////////////////////
	// send report
	jsonstr, jsonErr := json.Marshal(r)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonstr)
	}

	err := apiClient.Submit("POST", "/api/v1/pac/report/create", string(jsonstr), "report")
	if err != nil {
		log.Errorf("error: %+v\n", err)
		return err
	}

	////////////////////////////////////////////////////////////////////////
	// send namespaces
	jsonstr, jsonErr = json.Marshal(r.namespaces.Namespaces)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonErr)
	}

	err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/namespace/create", string(jsonstr), "namespace")
	if err != nil {
		log.Errorf("error: %+v\n", err)
		return err
	}

	////////////////////////////////////////////////////////////////////////
	// send pods
	jsonstr, jsonErr = json.Marshal(r.objectsList.pods)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonErr)
	}

	err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/pod/create", string(jsonstr), "pod")
	if err != nil {
		log.Errorf("error: %+v\n", err)
		return err
	}

	////////////////////////////////////////////////////////////////////////
	// send containers
	jsonstr, jsonErr = json.Marshal(r.objectsList.containers)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonErr)
	}

	err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/container/create", string(jsonstr), "container")
	if err != nil {
		log.Errorf("error: %+v\n", err)
		return err
	}

	////////////////////////////////////////////////////////////////////////
	// send containers
	jsonstr, jsonErr = json.Marshal(r.objectsList.uniqueImages)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonErr)
	}

	//err = apiClient.SendObjects(r.Uid, jsonstr)
	err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/image/create", string(jsonstr), "image")
	if err != nil {
		log.Errorf("error: %+v\n", err)
		return err
	}

	for _, images := range r.targetslist {
		////////////////////////////////////////////////////////////////////////
		// send targets
		for _, target := range images {
			var targetList []*Target // TODO Uggly hack to send a list of a single target
			targetList = append(targetList, target)

			jsonstr, jsonErr = json.Marshal(targetList)
			if jsonErr != nil {
				fmt.Printf("json error: %+v\n", jsonErr)
			}

			//err = apiClient.SendObjects(r.Uid, jsonstr)/api/v1/pac/report/{report_uid}/{image_uid}/vuln/create
			err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/"+target.ImageUid+"/vuln/create", string(jsonstr), "vuln")
			if err != nil {
				log.Errorf("error: %+v\n", err)
				return err
			}
		}
	}

	////////////////////////////////////////////////////////////////////////
	// send containers to image links
	jsonstr, jsonErr = json.Marshal(r.containerHasImage)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonErr)
	}

	//err = apiClient.SendObjects(r.Uid, jsonstr)
	err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/containerhasimage/create", string(jsonstr), "chi")
	if err != nil {
		log.Errorf("error: %+v\n", err)
		return err
	}

	for _, audit := range r.kubeauditReports {
		////////////////////////////////////////////////////////////////////////
		// send kubeaudit reports
		jsonstr, jsonErr = json.Marshal(audit)
		if jsonErr != nil {
			fmt.Printf("json error: %+v\n", jsonErr)
		}

		//err = apiClient.SendObjects(r.Uid, jsonstr)
		err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/audit/create", string(jsonstr), "audit")
		if err != nil {
			log.Errorf("error: %+v\n", err)
			return err
		}
	}

	////////////////////////////////////////////////////////////////////////
	// send containers to image links
	jsonstr, jsonErr = json.Marshal(r.reportSummary)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonErr)
	}

	//err = apiClient.SendObjects(r.Uid, jsonstr)
	err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/summary/create", string(jsonstr), "reportsummary")
	if err != nil {
		log.Errorf("error: %+v\n", err)
		return err
	}

	if opt.LimitDate > 0 || opt.LimitNr > 0 {
		var cleanup struct {
			LimitNr   int `json:"limit_nr"`
			LimitDate int `json:"limit_date"`
		}
		cleanup.LimitDate = opt.LimitDate
		cleanup.LimitNr = opt.LimitNr

		jsonstr, jsonErr = json.Marshal(cleanup)
		if jsonErr != nil {
			fmt.Printf("json error: %+v\n", jsonErr)
		}

		err = apiClient.Submit("POST", "/api/v1/pac/report/cleanup", string(jsonstr), "cleanup")
		if err != nil {
			log.Errorf("error: %+v\n", err)
			return err
		}
	}

	return nil

}

func (r *Report) Print() {

	color := map[string]string{
		"reset":        "\033[0m",
		"black":        "\033[30m",
		"red":          "\033[31m",
		"green":        "\033[32m",
		"yellow":       "\033[33m",
		"blue":         "\033[34m",
		"magenta":      "\033[35m",
		"cyan":         "\033[36m",
		"lightgray":    "\033[37m",
		"darkgray":     "\033[90m",
		"lightred":     "\033[91m",
		"lightgreen":   "\033[92m",
		"lightyellow":  "\033[93m",
		"lightblue":    "\033[94m",
		"lightmagenta": "\033[95m",
		"lightcyan":    "\033[96m",
		"white":        "\033[97m",

		"bgblack":        "\033[40m",
		"bgred":          "\033[41m",
		"bggreen":        "\033[42m",
		"bgyellow":       "\033[43m",
		"bgblue":         "\033[44m",
		"bgmagenta":      "\033[45m",
		"bgcyan":         "\033[46m",
		"bglightgray":    "\033[47m",
		"bgdarkgray":     "\033[100m",
		"bglightred":     "\033[101m",
		"bglightgreen":   "\033[102m",
		"bglightyellow":  "\033[103m",
		"bglightblue":    "\033[104m",
		"bglightmagenta": "\033[105m",
		"bglightcyan":    "\033[106m",
		"bgwhite":        "\033[107m",

		"bold": "\033[1m",
		"dim":  "\033[2m",
		"ul":   "\033[4m",
		"bl":   "\033[5m",
		"rev":  "\033[7m",
		"hid":  "\033[8m",
	}

	//fmt.Printf("Report %s\n", r.Uid)
	fmt.Printf("%sReport ===========================================%s\n", color["bold"], color["reset"])
	fmt.Printf("\tPods: %d\n", len(r.objectsList.pods))
	fmt.Printf("\tContainers: %d\n", len(r.objectsList.containers))
	fmt.Printf("\tImages: %d\n", len(r.objectsList.uniqueImages))
	fmt.Printf("\tTargets: %d\n", len(r.targetslist))

	for _, image := range r.objectsList.uniqueImages {
		fmt.Printf("\tImage: %s\n", image.Fulltag)
		fmt.Printf("\t\t%sTotal    : %d/%d%s\n", color["white"], image.Summary.Total, image.Summary.Fixed, color["reset"])
		fmt.Printf("\t\t%sCritical : %d/%d%s\n", color["red"], image.Summary.Severity.Critical.Total, image.Summary.Severity.Critical.Fixed, color["reset"])
		fmt.Printf("\t\t%sHigh     : %d/%d%s\n", color["yellow"], image.Summary.Severity.High.Total, image.Summary.Severity.High.Fixed, color["reset"])
		fmt.Printf("\t\t%sMedium   : %d/%d%s\n", color["cyan"], image.Summary.Severity.Medium.Total, image.Summary.Severity.Medium.Fixed, color["reset"])
		fmt.Printf("\t\t%sLow      : %d/%d%s\n", color["darkgray"], image.Summary.Severity.Low.Total, image.Summary.Severity.Low.Fixed, color["reset"])
		fmt.Printf("\t\t%sUnknown  : %d/%d%s\n", color["lightgray"], image.Summary.Severity.Unknown.Total, image.Summary.Severity.Unknown.Fixed, color["reset"])
		fmt.Println("")
	}

	fmt.Printf("%sReport Summary ===================================%s\n", color["bold"], color["reset"])
	fmt.Printf("\t       Total   : %v\n", r.reportSummary.VulnTotal)
	fmt.Printf("\t       Fixed   : %v\n", r.reportSummary.VulnFixed)
	fmt.Printf("\t       Critical: %d\n", r.reportSummary.VulnCritical)
	fmt.Printf("\t       High    : %d\n", r.reportSummary.VulnHigh)
	fmt.Printf("\t       Medium  : %d\n", r.reportSummary.VulnMedium)
	fmt.Printf("\t       Low     : %d\n", r.reportSummary.VulnLow)
	fmt.Printf("\t       Unknown : %d\n", r.reportSummary.VulnUnknown)
}
