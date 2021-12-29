package klustair

import (
	"encoding/json"
	"fmt"

	ka "github.com/Shopify/kubeaudit"
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/google/uuid"
	"github.com/klustair/klustair-go/pkg/api"
	"github.com/klustair/klustair-go/pkg/kubeaudit"
)

type Report struct {
	Uid              string `json:"uid"`
	Label            string `json:"label"`
	namespaces       *NamespaceList
	objectsList      *ObjectsList
	kubeauditReports []*ka.Report
	trivyreports     []*report.Report
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
		fmt.Printf("kubeaudit: %+v\n", len(kubeauditAuditors))
		k := new(kubeaudit.Auditor)
		nsList := r.namespaces.GetNamespaces()
		k.SetConfig(kubeauditAuditors)
		r.kubeauditReports = k.Run(nsList)
	}

	if trivy {
		uniqueImages := r.objectsList.GetUniqueImages()
		r.trivyreports, _ = Trivy.NewScanner().ScanImages(uniqueImages)
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

func (r *Report) SendReport(opt Options) error {
	fmt.Println("SEND REPORT ------------------")
	apiClient := api.NewApiClient(opt.Apihost, opt.Apitoken)

	fmt.Printf("namespaces: %+v\n", r.namespaces)
	reportJSON, jsonErr := json.Marshal(r.namespaces)
	if jsonErr != nil {
		fmt.Printf("json error: %+v\n", jsonErr)
	}

	err := apiClient.SendReport(reportJSON)
	if err != nil {
		return err
	}

	return nil

}
