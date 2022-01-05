package klustair

import (
	"encoding/json"
	"fmt"

	ka "github.com/Shopify/kubeaudit"
	"github.com/google/uuid"
	"github.com/klustair/klustair-go/pkg/api"
	"github.com/klustair/klustair-go/pkg/kubeaudit"
)

type Report struct {
	Uid              string `json:"uid"`
	Label            string `json:"title"`
	namespaces       *NamespaceList
	objectsList      *ObjectsList
	kubeauditReports []*ka.Report
	targetslist      Targetslist
	reportSummary    *ReportSummary
}

func (r *Report) Init(label string, whitelist []string, blacklist []string, trivy bool, kubeauditAuditors []string) {
	r.Uid = uuid.New().String()
	r.Label = label

	ns := new(NamespaceList)
	ns.Init(whitelist, blacklist)
	r.namespaces = ns

	o := new(ObjectsList)
	o.Init(r.namespaces)
	r.objectsList = o

	//kubeauditAuditors = nil
	if len(kubeauditAuditors) > 0 && kubeauditAuditors[0] != "" {
		k := new(kubeaudit.Auditor)
		nsList := r.namespaces.GetNamespaces()
		k.SetConfig(kubeauditAuditors)
		r.kubeauditReports = k.Run(nsList)
	}

	if trivy {
		r.targetslist = o.ScanImages()
	}

	rs := new(ReportSummary)
	rs.Init()
	rs.namespaces_total = ns.Total
	r.reportSummary = rs
}

func NewReport(opt Options) *Report {
	r := new(Report)
	r.Init(opt.Label, opt.Namespaces, opt.NamespacesBlacklist, opt.Trivy, opt.KubeAudit)
	return r
}

func (r *Report) Send(opt Options) error {
	fmt.Println("SEND REPORT ------------------")
	apiClient := api.NewApiClient(opt.Apihost, opt.Apitoken)

	////////////////////////////////////////////////////////////////////////////
	// send report
	jsonstr, jsonErr := json.Marshal(r)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonstr)
	}

	err := apiClient.Submit("POST", "/api/v1/pac/report/create", string(jsonstr), "report")
	if err != nil {
		return err
	}

	for _, namespace := range r.namespaces.Namespaces {
		////////////////////////////////////////////////////////////////////////
		// send namespaces
		jsonstr, jsonErr = json.Marshal(namespace)
		if jsonErr != nil {
			fmt.Printf("json error: %+v\n", jsonErr)
		}

		//err = apiClient.SendNamespaces(r.Uid, jsonstr)
		err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/namespace/create", string(jsonstr), "namespace")
		if err != nil {
			return err
		}
	}

	for _, pod := range r.objectsList.pods {
		////////////////////////////////////////////////////////////////////////
		// send pods
		jsonstr, jsonErr = json.Marshal(pod)
		if jsonErr != nil {
			fmt.Printf("json error: %+v\n", jsonErr)
		}

		//err = apiClient.SendObjects(r.Uid, jsonstr)
		err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/pod/create", string(jsonstr), "pod")
		if err != nil {
			return err
		}
	}

	for _, container := range r.objectsList.containers {
		////////////////////////////////////////////////////////////////////////
		// send containers
		jsonstr, jsonErr = json.Marshal(container)
		if jsonErr != nil {
			fmt.Printf("json error: %+v\n", jsonErr)
		}

		//err = apiClient.SendObjects(r.Uid, jsonstr)
		err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/container/create", string(jsonstr), "container")
		if err != nil {
			return err
		}
	}

	for _, image := range r.objectsList.uniqueImages {
		////////////////////////////////////////////////////////////////////////
		// send containers
		jsonstr, jsonErr = json.Marshal(image)
		if jsonErr != nil {
			fmt.Printf("json error: %+v\n", jsonErr)
		}

		//err = apiClient.SendObjects(r.Uid, jsonstr)
		err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/image/create", string(jsonstr), "image")
		if err != nil {
			return err
		}
	}

	for image, target := range r.targetslist {
		////////////////////////////////////////////////////////////////////////
		// send containers
		jsonstr, jsonErr = json.Marshal(target)
		if jsonErr != nil {
			fmt.Printf("json error: %+v\n", jsonErr)
		}

		//err = apiClient.SendObjects(r.Uid, jsonstr)/api/v1/pac/report/{report_uid}/{image_uid}/vuln/create
		err = apiClient.Submit("POST", "/api/v1/pac/report/"+r.Uid+"/"+image+"/target/create", string(jsonstr), "target")
		if err != nil {
			return err
		}
	}

	return nil

}
