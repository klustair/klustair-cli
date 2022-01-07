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

	// run kubeaudit scans if enabled
	if len(kubeauditAuditors) > 0 && kubeauditAuditors[0] != "" {
		k := new(kubeaudit.Auditor)
		k.Klustair.ReportUid = r.Uid
		k.SetConfig(kubeauditAuditors)

		for _, namespace := range r.namespaces.Namespaces {
			log.Debugf("Kubeaudit on namespace: %+v", namespace.Name)
			report := k.Audit(namespace.Name)
			k.Klustair.NamespaceUid = namespace.Uid

			r.kubeauditReports = append(r.kubeauditReports, report)
		}
	}

	// run trivy scans if enabled
	if trivy {
		r.targetslist = o.ScanImages()
	}

	r.containerHasImage = o.linkImagesToContainers(trivy)

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
	return nil

}
