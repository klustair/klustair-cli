package klustair

import "github.com/google/uuid"

type Target struct {
	Uid             string           `json:"uid"`
	ReportUid       string           `json:"report_uid"`
	ImageUid        string           `json:"image_uid"`
	Target          string           `json:"Target"`
	TargetType      string           `json:"Type"`
	IsOS            bool             `json:"isOS"`
	Vulnerabilities []*Vulnerability `json:"Vulnerabilities"`
	summary         VulnSummary
}

func NewTarget(reportUid string, imageUid string) *Target {
	t := new(Target)
	t.Uid = uuid.New().String()
	t.ReportUid = reportUid
	t.ImageUid = imageUid
	return t
}
