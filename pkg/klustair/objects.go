package klustair

import (
	"fmt"

	"github.com/aquasecurity/trivy/pkg/report"
)

type ObjectsList struct {
	pods         []*Pod
	containers   []*Container
	images       []*Image
	trivyreports []*report.Report
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

func (ol *ObjectsList) GetUniqueImages() map[string]*Image {
	//var unique map[]images
	uniqueImages := make(map[string]*Image)

	for _, image := range ol.images {
		uniqueImages[image.fulltag] = image
	}
	return uniqueImages
}

func (ol *ObjectsList) ScanImages() {
	//var unique map[]images
	uniqueImages := ol.GetUniqueImages()
	for _, image := range uniqueImages {
		fmt.Println("fulltag:", image.fulltag)
		report, err := image.Scan()
		if err != nil {
			fmt.Errorf("error scanning image: %s", err)
		}
		ol.trivyreports = append(ol.trivyreports, &report)
	}

}
