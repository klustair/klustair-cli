package klustair

import (
	"github.com/google/uuid"
)

type Report struct {
	uid           string
	label         string
	namespaces    *NamespaceList
	objectsList   *ObjectsList
	reportSummary *ReportSummary
}

func (r *Report) Init(label string, whitelist []string, blacklist []string, trivy bool) {
	r.uid = uuid.New().String()
	r.label = label

	ns := new(NamespaceList)
	ns.Init(whitelist, blacklist)
	r.namespaces = ns

	o := new(ObjectsList)
	o.Init(r.namespaces)
	r.objectsList = o

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
	r.Init(opt.Label, opt.Namespaces, opt.NamespacesBlacklist, opt.Trivy)
	return r
}
