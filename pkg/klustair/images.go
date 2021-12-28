package klustair

import (
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/google/uuid"
	"github.com/klustair/klustair-go/pkg/trivyscanner"
)

var Trivy *trivyscanner.Trivy

type Image struct {
	uid            string
	fulltag        string
	image_b64      string
	arch           string
	layer_count    int
	image_digest   string
	distro         string
	distro_version string
	age            int
	config         string
	history        string
}

func (i *Image) Init(fulltag string) {
	i.uid = uuid.New().String()
	i.fulltag = fulltag
}

func (i *Image) Scan() (report.Report, error) {
	trivy := Trivy.NewScanner()
	return trivy.Scan(i.fulltag)
}
