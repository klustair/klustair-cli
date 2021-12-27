package klustair

import (
	"github.com/google/uuid"
)

type Report struct {
	uid           string
	label         string
	namespaces    *NamespaceList
	podsList      *PodsList
	reportSummary *ReportSummary

	//clientset *kubernetes.Clientset
}

func (r *Report) Init(label string, whitelist []string, blacklist []string) {
	r.uid = uuid.New().String()
	r.label = label

	ns := new(NamespaceList)
	ns.Init(whitelist, blacklist)
	r.namespaces = ns

	p := new(PodsList)
	p.Init(r.namespaces)
	r.podsList = p

	rs := new(ReportSummary)
	rs.Init()
	rs.namespaces_total = ns.total
	r.reportSummary = rs
}

func NewReport(opt Options) *Report {
	r := new(Report)
	r.Init(opt.Label, opt.Namespaces, opt.NamespacesBlacklist)
	return r
}
