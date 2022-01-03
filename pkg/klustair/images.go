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
}

func (i *Image) Init(fulltag string) {
	i.Uid = uuid.New().String()
	i.Fulltag = fulltag
	log.Debugf("    image: %+s", fulltag)
}

func (i *Image) Scan() (report.Report, error) {
	trivy := Trivy.NewScanner()
	// set current timestamp
	i.AnalyzedAt = time.Now().Unix()

	// scan image
	report, err := trivy.Scan(i.Fulltag)
	if err != nil {
		log.Errorf("    trivy failed to scan image %+v: %+v", i.Fulltag, err)
	}
	i.LayerCount = len(report.Metadata.ImageConfig.RootFS.DiffIDs)
	return report, err
}
