package klustair

import (
	"github.com/google/uuid"
)

type ReportSummary struct {
	uid                string
	namespaces_total   int
	namespaces_checked int
	vuln_total         int
	vuln_high          int
	vuln_critical      int
	vuln_medium        int
	vuln_low           int
	vuln_unknown       int
	vuln_negligible    int
	vuln_fixed         int
	pods               int
	containers         int
	images             int
}

func (rs *ReportSummary) Init() {
	rs.uid = uuid.New().String()
	rs.namespaces_total = 0
	rs.namespaces_checked = 0
	rs.vuln_total = 0
	rs.vuln_high = 0
	rs.vuln_critical = 0
	rs.vuln_medium = 0
	rs.vuln_low = 0
	rs.vuln_unknown = 0
	rs.vuln_negligible = 0
	rs.vuln_fixed = 0
	rs.pods = 0
	rs.containers = 0
	rs.images = 0
}

func (rs *ReportSummary) sumVulnSummary(uniqueImages map[string]*Image) {
	for _, image := range uniqueImages {
		rs.vuln_total += image.summary.total
		rs.vuln_high += image.summary.high
		rs.vuln_critical += image.summary.critical
		rs.vuln_medium += image.summary.medium
		rs.vuln_low += image.summary.low
		rs.vuln_unknown += image.summary.unknown
		rs.vuln_fixed += image.summary.fixed
	}
}

func NewReportSummary() *ReportSummary {
	r := new(ReportSummary)
	r.Init()
	return r
}
