package klustair

import (
	"fmt"

	ka "github.com/Shopify/kubeaudit"
	"github.com/google/uuid"
	"github.com/klustair/klustair-go/pkg/kubeaudit"
)

type Report struct {
	uid             string
	label           string
	namespaces      *NamespaceList
	objectsList     *ObjectsList
	kubeauditReport *ka.Report
	reportSummary   *ReportSummary
}

func (r *Report) Init(label string, whitelist []string, blacklist []string, trivy bool, kubeauditAuditors string) {
	r.uid = uuid.New().String()
	r.label = label

	ns := new(NamespaceList)
	ns.Init(whitelist, blacklist)
	r.namespaces = ns

	o := new(ObjectsList)
	o.Init(r.namespaces)
	r.objectsList = o

	//kubeauditAuditors = nil
	if kubeauditAuditors != "" {
		fmt.Printf("kubeaudit: %+v\n", kubeauditAuditors)
		k := new(kubeaudit.Auditor)
		k.SetConfig(kubeauditAuditors)
		r.kubeauditReport = k.Run()
	}

	if trivy {
		o.ScanImages()
	}

	rs := new(ReportSummary)
	rs.Init()
	rs.namespaces_total = ns.total
	r.reportSummary = rs
}

func NewReport(opt Options) *Report {
	r := new(Report)
	r.Init(opt.Label, opt.Namespaces, opt.NamespacesBlacklist, opt.Trivy, opt.KubeAudit)
	return r
}
