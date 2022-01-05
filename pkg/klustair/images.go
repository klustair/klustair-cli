package klustair

import (
	"time"

	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/google/uuid"
	"github.com/klustair/klustair-go/pkg/trivyscanner"
	log "github.com/sirupsen/logrus"
)

var Trivy *trivyscanner.Trivy

type Image struct {
	Uid           string `json:"uid"`
	ReportUid     string `json:"report_uid"`
	Image_b64     string `json:"image_b64"`
	AnalyzedAt    int64  `json:"analyzed_at"`
	CreatedAt     int64  `json:"created_at"`
	Fulltag       string `json:"fulltag"`
	ImageDigest   string `json:"image_digest"`
	Arch          string `json:"arch"`
	Distro        string `json:"distro"`
	DistroVersion string `json:"distro_version"`
	ImageSize     int    `json:"image_size"`
	LayerCount    int    `json:"layer_count"`
	Registry      string `json:"registry"`
	Repo          string `json:"repo"`
	Dockerfile    string `json:"dockerfile"`
	Config        string `json:"config"`
	History       string `json:"history"`
	Age           int    `json:"age"`
	Targets       []trivyscanner.Target
}

func (i *Image) Init(fulltag string) {
	i.Uid = uuid.New().String()
	i.Fulltag = fulltag
	log.Debugf("    image: %+s", fulltag)
}

func (i *Image) Scan() ([]trivyscanner.Target, error) {
	trivy := Trivy.NewScanner()

	// scan image
	report, err := trivy.Scan(i.Fulltag)
	if err != nil {
		log.Errorf("    trivy failed to scan image %+v: %+v", i.Fulltag, err)
		return nil, err
	}
	i.Arch = report.Metadata.ImageConfig.Architecture
	i.LayerCount = len(report.Metadata.ImageConfig.RootFS.DiffIDs)
	i.ImageDigest = report.Metadata.RepoDigests[0]
	i.Distro = report.Metadata.OS.Family
	i.DistroVersion = report.Metadata.OS.Name
	i.CreatedAt = report.Metadata.ImageConfig.Created.Unix()
	i.AnalyzedAt = time.Now().Unix()
	i.Age = int(time.Now().Sub(time.Unix(i.CreatedAt, 0)).Hours() / 24)
	// TODO Find a way to save those informations
	//i.Config = report.Metadata.ImageConfig.Config
	//i.History = report.Metadata.ImageConfig.History
	targets := i.getVulnerabilities(report)
	return targets, err
}

func (i *Image) getVulnerabilities(report report.Report) []trivyscanner.Target {
	var targets []trivyscanner.Target
	for _, target := range report.Results {
		t := trivyscanner.Target{
			Vulnerabilities: []trivyscanner.Vulnerability{},
		}
		for _, vuln := range target.Vulnerabilities {
			//TODO delete me
			//fmt.Printf("CVSS:%+v\n", vuln.CVSS)
			v := trivyscanner.Vulnerability{
				VulnerabilityID:  vuln.VulnerabilityID,
				PkgName:          vuln.PkgName,
				InstalledVersion: vuln.InstalledVersion,
				FixedVersion:     vuln.FixedVersion,
				Title:            vuln.Title,
				Description:      vuln.Description,
				Severity:         vuln.Severity,
				SeveritySource:   vuln.SeveritySource,
				LastModifiedDate: vuln.LastModifiedDate,
				PublishedDate:    vuln.PublishedDate,
				References:       vuln.References,
				// TODO fill with cvss
				//CVSS:             vuln.CVSS,
				CweIDs: vuln.CweIDs,
			}
			t.Vulnerabilities = append(t.Vulnerabilities, v)
		}
		//i.Targets = append(i.Targets, t)
		targets = append(targets, t)
	}
	return targets
}
