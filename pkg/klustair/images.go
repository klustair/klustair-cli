package klustair

import (
	"fmt"
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
	AnalyzedAt    string `json:"analyzed_at"`
	CreatedAt     string `json:"created_at"`
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
	//Targets       []*Target
	Summary VulnSummary `json:"summary"`
}

func (i *Image) Init(fulltag string) {
	i.Uid = uuid.New().String()
	i.Fulltag = fulltag
	i.Config = "{}"
	i.History = "{}"
	log.Debugf("    image: %+s", fulltag)
}

func (i *Image) Scan() ([]*Target, error) {
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
	i.CreatedAt = report.Metadata.ImageConfig.Created.UTC().Format(time.RFC3339)
	i.AnalyzedAt = time.Now().UTC().Format(time.RFC3339)
	i.Age = int(time.Now().Sub(time.Unix(report.Metadata.ImageConfig.Created.Unix(), 0)).Hours() / 24)
	// TODO Find a way to save those informations
	//i.Config = report.Metadata.ImageConfig.Config
	//i.History = report.Metadata.ImageConfig.History
	//i.Dockerfile = report.Metadata.....
	//i.Repo = report.Metadata.RepoName
	//i.Registry = report.Metadata.RepoName
	targets := i.getVulnerabilities(report)
	return targets, err
}

func (i *Image) getVulnerabilities(report report.Report) []*Target {
	var targets []*Target
	for _, target := range report.Results {
		t := NewTarget(i.ReportUid, i.Uid)
		t.Target = target.Target
		t.TargetType = target.Type
		for _, vuln := range target.Vulnerabilities {
			//TODO delete me
			//fmt.Printf("CVSS:%+v\n", vuln.CVSS)
			v := NewVulnerability(i.ReportUid, i.Uid, t.Uid)
			v.VulnerabilityID = vuln.VulnerabilityID
			v.PkgName = vuln.PkgName
			v.Title = vuln.Title
			v.Description = vuln.Description
			v.InstalledVersion = vuln.InstalledVersion
			v.FixedVersion = vuln.FixedVersion
			v.SeveritySource = vuln.SeveritySource
			v.Severity = vuln.Severity
			v.LastModifiedDate = vuln.LastModifiedDate
			v.PublishedDate = vuln.PublishedDate
			v.References = vuln.References
			// TODO fill with cvss
			//v.CVSS = vuln.CVSS
			v.CweIDs = vuln.CweIDs

			// Increment summary
			i.Summary.Add(v)

			t.Vulnerabilities = append(t.Vulnerabilities, v)
		}
		//i.Targets = append(i.Targets, t)
		targets = append(targets, t)
	}
	fmt.Printf("summary    %+v\n", i.Summary)
	return targets
}
