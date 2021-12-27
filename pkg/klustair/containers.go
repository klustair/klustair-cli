package klustair

import (
	"fmt"

	"github.com/google/uuid"
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

	for _, status := range containerstatus {
		if status.Name == c.name {
			fmt.Printf("STATUS: %+v\n", status.Name)
			c.ready = status.Ready
			//c.started = status.State.Running != nil
			c.restartCount = status.RestartCount
			c.imageID = status.ImageID
			//c.startedAt = status.State.Running.StartedAt.Unix()
			c.actual = true
		}
	}
	fmt.Println("container: ", c.name, c.ready)
}
