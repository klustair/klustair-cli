package klustair

import (
	"github.com/google/uuid"
)

type ReportSummary struct {
	uid               string
	NamespacesTotal   int `json:"namespaces_total"`
	NamespacesChecked int `json:"namespaces_checked"`
	VulnTotal         int `json:"vuln_total"`
	VulnHigh          int `json:"vuln_high"`
	VulnCritical      int `json:"vuln_critical"`
	VulnMedium        int `json:"vuln_medium"`
	VulnLow           int `json:"vuln_low"`
	VulnUnknown       int `json:"vuln_unknown"`
	VulnNegligible    int `json:"vuln_negligible"`
	VulnFixed         int `json:"vuln_fixed"`
	Pods              int `json:"pods"`
	Containers        int `json:"containers"`
	Images            int `json:"images"`
}

func (rs *ReportSummary) Init() {
	rs.uid = uuid.New().String()
	rs.NamespacesTotal = 0
	rs.NamespacesChecked = 0
	rs.VulnTotal = 0
	rs.VulnFixed = 0
	rs.VulnCritical = 0
	rs.VulnHigh = 0
	rs.VulnMedium = 0
	rs.VulnLow = 0
	rs.VulnUnknown = 0
	rs.VulnNegligible = 0
	rs.Pods = 0
	rs.Containers = 0
	rs.Images = 0
}

func (rs *ReportSummary) sumVulnSummary(uniqueImages map[string]*Image) {
	for _, image := range uniqueImages {
		s := image.Summary
		rs.VulnTotal += s.Total
		rs.VulnFixed += s.Fixed
		rs.VulnCritical += s.Severity.Critical.Total
		rs.VulnHigh += s.Severity.High.Total
		rs.VulnMedium += s.Severity.Medium.Total
		rs.VulnLow += s.Severity.Low.Total
		rs.VulnUnknown += s.Severity.Unknown.Total
	}
}

func NewReportSummary() *ReportSummary {
	r := new(ReportSummary)
	r.Init()
	return r
}
