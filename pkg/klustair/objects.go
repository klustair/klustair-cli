package klustair

import (
	log "github.com/sirupsen/logrus"
)

type ObjectsList struct {
	pods         []*Pod
	containers   []*Container
	uniqueImages map[string]*Image
}

type Targetslist map[string][]*Target

func (ol *ObjectsList) Init(reportUid string, namespaces *NamespaceList) {
	ol.uniqueImages = make(map[string]*Image)

	for _, namespace := range namespaces.Namespaces {

		podsList, err := Client.GetPods(namespace.Name)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, pod := range podsList.Items {

			p := new(Pod)
			p.Init(reportUid, namespace.Uid, pod)

			log.Debug("pod:", p.Podname)
			ol.pods = append(ol.pods, p)

			for _, container := range pod.Spec.Containers {
				c := new(Container)
				c.Init(reportUid, namespace.Uid, p.Uid, container, pod.Status.ContainerStatuses, false)
				//fmt.Printf("container: %+v\n", c)
				ol.containers = append(ol.containers, c)

				i := new(Image)
				i.Init(container.Image)
				//ol.images = append(ol.images, i)
				ol.uniqueImages[container.Image] = i
			}

			for _, initcontainer := range pod.Spec.InitContainers {
				c := new(Container)
				c.Init(reportUid, namespace.Uid, p.Uid, initcontainer, pod.Status.ContainerStatuses, true)
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

func (ol *ObjectsList) ScanImages() Targetslist {
	//var unique map[]images
	trivyReports := make(Targetslist)
	for _, image := range ol.uniqueImages {
		log.Info("Trivy scan image:", image.Fulltag)
		targets, err := image.Scan()
		if err != nil {
			log.Errorf("error scanning image: %s", err)
		}
		trivyReports[image.Uid] = targets
	}
	return trivyReports
}

type ContainerHasImage struct {
	ReportUid    string `json:"report_uid"`
	ContainerUid string `json:"container_uid"`
	ImageUid     string `json:"image_uid"`
}

func (ol *ObjectsList) linkImagesToContainers(trivyScan bool) []*ContainerHasImage {
	var containerHasImages []*ContainerHasImage

	// iterate over all container and find the corresponding image
	for _, container := range ol.containers {
		for _, image := range ol.uniqueImages {
			if container.Image == image.Fulltag {
				containerHasImages = append(containerHasImages, &ContainerHasImage{
					ReportUid:    container.ReportUid,
					ContainerUid: container.Uid,
					ImageUid:     image.Uid,
				})
			}

			// works only after a trivy scan
			if trivyScan {
				// TODO needs some refinment sice they are not completly equal
				if container.ImageID == image.ImageDigest {
					container.Actual = true
				} else {
					container.Actual = false
				}
			}

		}
	}
	return containerHasImages
}
