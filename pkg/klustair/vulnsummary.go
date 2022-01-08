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
	total    int
	fixed    int
	critical int
	high     int
	medium   int
	low      int
	unknown  int
}

func (v *VulnSummary) Add(vulnerability *Vulnerability) {
	v.total++
	switch vulnerability.Severity {
	case SeverityCritical:
		v.critical++
	case SeverityHigh:
		v.high++
	case SeverityMedium:
		v.medium++
	case SeverityLow:
		v.low++
	case SeverityUnknown:
		v.unknown++
	}

	if vulnerability.FixedVersion != "" {
		v.fixed++
	}
}
