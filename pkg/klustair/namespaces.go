package klustair

import (
	"fmt"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
)

type Namespace struct {
	Name                     string    `json:"name"`
	Uid                      string    `json:"uid"`
	Kubernetes_namespace_uid types.UID `json:"kubernetes_namespace_uid"`
	Creation_timestamp       int64     `json:"creation_timestamp"`
}

type NamespaceList struct {
	Namespaces []Namespace `json:"namespaces"`
	Total      int         `json:"total"`
	Checked    int         `json:"checked"`
}

func (n *Namespace) Init(name string, kubernetes_namespace_uid types.UID, creation_timestamp int64) {
	n.Name = name
	n.Uid = uuid.New().String()
	n.Kubernetes_namespace_uid = kubernetes_namespace_uid
	n.Creation_timestamp = creation_timestamp
}

func (ns *NamespaceList) Init(whitelist []string, blacklist []string) {

	namespaceList, err := Client.GetNamespaces()
	if err != nil {
		panic(err)
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
		n.Init(namespace.Name, namespace.UID, namespace.CreationTimestamp.Unix())

		// TODO remove me
		fmt.Printf("namespace: %+v\n", n)
		ns.Namespaces = append(ns.Namespaces, *n)
	}
	ns.Checked = len(ns.Namespaces)

}

func (ns *NamespaceList) GetNamespaces() []string {
	var namespacesList []string
	for _, namespace := range ns.Namespaces {
		fmt.Println("namespace:", namespace.Name)
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
