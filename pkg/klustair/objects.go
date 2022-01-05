package klustair

import (
	"os"

	"github.com/klustair/klustair-go/pkg/trivyscanner"
	log "github.com/sirupsen/logrus"
)

type ObjectsList struct {
	pods         []*Pod
	containers   []*Container
	uniqueImages map[string]*Image
}

type Targetslist map[string][]trivyscanner.Target

func (ol *ObjectsList) Init(namespaces *NamespaceList) {
	ol.uniqueImages = make(map[string]*Image)

	for _, namespace := range namespaces.Namespaces {

		podsList, err := Client.GetPods(namespace.Name)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		for _, pod := range podsList.Items {

			p := new(Pod)
			p.Init(namespace.Uid, pod)

			// TODO remove me
			//fmt.Printf("pod: %+v\n", p)
			log.Debug("pod:", p.Podname)
			ol.pods = append(ol.pods, p)

			for _, container := range pod.Spec.Containers {
				c := new(Container)
				c.Init(container, pod.Status.ContainerStatuses, false)
				//fmt.Printf("container: %+v\n", c)
				ol.containers = append(ol.containers, c)

				i := new(Image)
				i.Init(container.Image)
				//ol.images = append(ol.images, i)
				ol.uniqueImages[container.Image] = i
			}

			for _, initcontainer := range pod.Spec.InitContainers {
				c := new(Container)
				c.Init(initcontainer, pod.Status.ContainerStatuses, true)
				//fmt.Printf("initcontainer: %+v\n", c)
				ol.containers = append(ol.containers, c)

				i := new(Image)
				i.Init(initcontainer.Image)
				//ol.images = append(ol.images, i)
				ol.uniqueImages[initcontainer.Image] = i
			}

		}
	}
}

func (ol *ObjectsList) ScanImages() Targetslist { //replace String with trivy report object
	//var unique map[]images
	trivyReports := make(Targetslist)
	for _, image := range ol.uniqueImages {
		log.Debug("Trivy scan image fulltag:", image.Fulltag)
		targets, err := image.Scan()
		if err != nil {
			log.Errorf("error scanning image: %s", err)
		}
		trivyReports[image.Uid] = targets
	}
	return trivyReports
}
