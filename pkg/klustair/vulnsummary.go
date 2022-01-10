package klustair

/*
type ImageSummary struct {
	Uid       string `json:"uid"`
	ReportUid string `json:"report_uid"`
	ImageUid  string `json:"image_uid"`
	Severity  string `json:"severity"`
	Total     int    `json:"total"`
	Fixed     int    `json:"fixed"`
}
*/

type VulnSummary struct {
	Total    int `json:"Total"`
	Fixed    int `json:"Fixed"`
	Severity struct {
		Critical Severity `json:"Critical"`
		High     Severity `json:"High"`
		Medium   Severity `json:"Medium"`
		Low      Severity `json:"Low"`
		Unknown  Severity `json:"Unknown"`
	} `json:"severity"`
}

type Severity struct {
	Total int `json:"total"`
	Fixed int `json:"fixed"`
}

func (v *VulnSummary) Add(vulnerability *Vulnerability) {
	v.Total++
	if vulnerability.FixedVersion != "" {
		v.Fixed++
	}

	switch vulnerability.Severity {
	case SeverityCritical:
		v.Severity.Critical.Total++
		if vulnerability.FixedVersion != "" {
			v.Severity.Critical.Fixed++
		}
	case SeverityHigh:
		v.Severity.High.Total++
		if vulnerability.FixedVersion != "" {
			v.Severity.High.Fixed++
		}
	case SeverityMedium:
		v.Severity.Medium.Total++
		if vulnerability.FixedVersion != "" {
			v.Severity.Medium.Fixed++
		}
	case SeverityLow:
		v.Severity.Low.Total++
		if vulnerability.FixedVersion != "" {
			v.Severity.Low.Fixed++
		}
	case SeverityUnknown:
		v.Severity.Unknown.Total++
		if vulnerability.FixedVersion != "" {
			v.Severity.Unknown.Fixed++
		}
	}
}
