package klustair

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
)

type Namespace struct {
	Name                     string    `json:"name"`
	Uid                      string    `json:"uid"`
	Kubernetes_namespace_uid types.UID `json:"kubernetes_namespace_uid"`
	Creation_timestamp       string    `json:"creation_timestamp"`
}

type NamespaceList struct {
	Namespaces []Namespace `json:"namespaces"`
	Total      int         `json:"total"`
	Checked    int         `json:"checked"`
}

func (n *Namespace) Init(name string, kubernetes_namespace_uid types.UID, creation_timestamp string) {
	n.Name = name
	n.Uid = uuid.New().String()
	n.Kubernetes_namespace_uid = kubernetes_namespace_uid
	n.Creation_timestamp = creation_timestamp
}

func (ns *NamespaceList) Init(whitelist []string, blacklist []string) {

	namespaceList, err := Client.GetNamespaces()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ns.Total = len(namespaceList.Items)
	for _, namespace := range namespaceList.Items {

		if len(whitelist) > 0 {
			if !stringInSlice(namespace.Name, whitelist) {
				continue
			}
		}

		if len(blacklist) > 0 {
			if stringInSlice(namespace.Name, blacklist) {
				continue
			}
		}

		n := new(Namespace)
		n.Init(namespace.Name, namespace.UID, namespace.CreationTimestamp.UTC().Format(time.RFC3339))

		log.Debug("namespace:", n.Name)
		ns.Namespaces = append(ns.Namespaces, *n)
	}
	ns.Checked = len(ns.Namespaces)

}

func (ns *NamespaceList) GetNamespaces() []string {
	var namespacesList []string
	for _, namespace := range ns.Namespaces {
		namespacesList = append(namespacesList, namespace.Name)
	}
	return namespacesList
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
