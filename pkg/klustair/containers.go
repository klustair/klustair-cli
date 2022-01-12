package klustair

import (
	"encoding/json"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type Container struct {
	Uid               string `json:"uid"`
	ReportUid         string `json:"report_uid"`
	NamespaceUid      string `json:"namespace_uid"`
	Pod_uid           string `json:"pod_uid"`
	Name              string `json:"name"`
	Image             string `json:"image"`
	Image_pull_policy string `json:"image_pull_policy"`
	Security_context  string `json:"security_context"`
	Init_container    bool   `json:"init_container"`
	Ready             bool   `json:"ready"`
	Started           bool   `json:"started"`
	RestartCount      int32  `json:"restartCount"`
	ImageID           string `json:"imageID"`
	StartedAt         int64  `json:"startedAt"`
	Actual            bool   `json:"actual"`
}

/*type ContainersList struct {
	containers []Container
}
*/

func (c *Container) Init(reportUid string, namespaceUid string, podUid string, container v1.Container, containerstatus []v1.ContainerStatus, init_container bool) {

	sc, scErr := json.Marshal(container.SecurityContext)
	if scErr != nil {
		log.Warnf("Error marshalling security context of %s: %s", container.Name, scErr)
	}

	c.Uid = uuid.New().String()
	c.ReportUid = reportUid
	c.NamespaceUid = namespaceUid
	c.Pod_uid = podUid
	c.Name = container.Name
	c.Image = container.Image
	c.Image_pull_policy = string(container.ImagePullPolicy)
	c.Security_context = string(sc)
	c.Init_container = init_container

	// TODO: This part needs some refinement (Missing fields and unusual values)
	for _, status := range containerstatus {
		if status.Name == c.Name {
			c.Ready = status.Ready
			//c.started = status.State.Running != nil
			c.RestartCount = status.RestartCount
			c.ImageID = status.ImageID
			//c.startedAt = status.State.Running.StartedAt.Unix()
			c.Actual = true
		}
	}
	log.Debugf("  container: %+s, ready: %+v", c.Name, c.Ready)
}
