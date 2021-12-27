package klustair

import (
	"fmt"

	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"
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

type ObjectsList struct {
	pods       []*Pod
	containers []*Container
	images     []*Image
}

func (p *Pod) Init(namespace_uid string, pod v1.Pod) {
	p.podname = pod.Name
	p.uid = uuid.New().String()
	p.namespace_uid = namespace_uid
	p.kubernetes_pod_uid = pod.UID
	p.creation_timestamp = pod.CreationTimestamp.Unix()
}

func (ol *ObjectsList) Init(namespaces *NamespaceList) {

	for _, namespace := range namespaces.namespaces {

		podsList, err := Client.GetPods(namespace.name)
		if err != nil {
			panic(err)
		}

		for _, pod := range podsList.Items {

			p := new(Pod)
			p.Init(namespace.uid, pod)

			// TODO remove me
			//fmt.Printf("pod: %+v\n", p)
			fmt.Println("pod:", p.podname)
			ol.pods = append(ol.pods, p)

			for _, container := range pod.Spec.Containers {
				c := new(Container)
				c.Init(container, pod.Status.ContainerStatuses, false)
				//fmt.Printf("container: %+v\n", c)
				ol.containers = append(ol.containers, c)

				i := new(Image)
				i.Init(container.Image)
				ol.images = append(ol.images, i)
			}

			for _, initcontainer := range pod.Spec.InitContainers {
				c := new(Container)
				c.Init(initcontainer, pod.Status.ContainerStatuses, true)
				//fmt.Printf("initcontainer: %+v\n", c)
				ol.containers = append(ol.containers, c)

				i := new(Image)
				i.Init(initcontainer.Image)
				ol.images = append(ol.images, i)
			}

		}
	}
}
