package klustair

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type Container struct {
	uid               string
	name              string
	report_uid        string
	namespace_uid     string
	pod_uid           string
	image             string
	image_pull_policy string
	//security_context  json.RawMessage
	init_container bool
	ready          bool
	started        bool
	restartCount   int32
	imageID        string
	startedAt      int64
	actual         bool
}

type ContainersList struct {
	containers []Container
}

func (c *Container) Init(container v1.Container, containerstatus []v1.ContainerStatus, init_container bool) {
	c.uid = uuid.New().String()
	c.name = container.Name
	c.image = container.Image
	c.image_pull_policy = string(container.ImagePullPolicy)
	//c.security_context = json.Unmarshal(container.SecurityContext)
	c.init_container = init_container

	// TODO: This part needs some refinement (Missing fields and unusual values)
	for _, status := range containerstatus {
		if status.Name == c.name {
			c.ready = status.Ready
			//c.started = status.State.Running != nil
			c.restartCount = status.RestartCount
			c.imageID = status.ImageID
			//c.startedAt = status.State.Running.StartedAt.Unix()
			c.actual = true
		}
	}
	log.Debugf("  container: %+s, ready: %+v", c.name, c.ready)
}
