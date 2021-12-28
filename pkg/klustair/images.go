package klustair

import (
	"fmt"

	"github.com/google/uuid"
)

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

func (i *Image) Scan() {
	trivy := Trivy.NewScanner()
	report, err := trivy.Scan(i.fulltag)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%+v\n", report)
}
