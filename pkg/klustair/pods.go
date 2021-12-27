package klustair

import (
	"fmt"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
)

type Pod struct {
	podname            string
	uid                string
	namespace_uid      string
	kubernetes_pod_uid types.UID
	creation_timestamp int64
	age                int
}

type PodsList struct {
	pods []Pod
	//containers map[string]Container
	containers []*Container
}

func (p *Pod) Init(name string, namespace_uid string, kubernetes_pod_uid types.UID, creation_timestamp int64) {
	p.podname = name
	p.uid = uuid.New().String()
	p.namespace_uid = namespace_uid
	p.kubernetes_pod_uid = kubernetes_pod_uid
	p.creation_timestamp = creation_timestamp
}

func (pl *PodsList) Init(namespaces *NamespaceList) {

	for _, namespace := range namespaces.namespaces {

		podsList, err := Client.GetPods(namespace.name)
		if err != nil {
			panic(err)
		}

		for _, pod := range podsList.Items {

			p := new(Pod)
			p.Init(pod.Name, namespace.uid, pod.UID, pod.CreationTimestamp.Unix())

			// TODO remove me
			//fmt.Printf("pod: %+v\n", p)
			fmt.Println("pod:", p.podname)
			pl.pods = append(pl.pods, *p)

			for _, container := range pod.Spec.Containers {
				c := new(Container)
				c.Init(container, pod.Status.ContainerStatuses, false)

				//fmt.Printf("container: %+v\n", c)
				pl.containers = append(pl.containers, c)
			}

			for _, initcontainer := range pod.Spec.InitContainers {
				c := new(Container)
				c.Init(initcontainer, pod.Status.ContainerStatuses, true)

				//fmt.Printf("initcontainer: %+v\n", c)
				pl.containers = append(pl.containers, c)
			}

		}
	}
}
