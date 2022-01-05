package trivyscanner

type Target struct {
	Uid             string          `json:"uid"`
	ReportUid       string          `json:"report_uid"`
	ImageUid        string          `json:"image_uid"`
	Target          string          `json:"Target"`
	TargetType      string          `json:"TargetType"`
	IsOS            bool            `json:"IsOS"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}
