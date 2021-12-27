package klustair

import (
	"fmt"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
)

type Namespace struct {
	name                     string
	uid                      string
	kubernetes_namespace_uid types.UID
	creation_timestamp       int64
}

type NamespaceList struct {
	namespaces []Namespace
	total      int
	checked    int
}

func (n *Namespace) Init(name string, kubernetes_namespace_uid types.UID, creation_timestamp int64) {
	n.name = name
	n.uid = uuid.New().String()
	n.kubernetes_namespace_uid = kubernetes_namespace_uid
	n.creation_timestamp = creation_timestamp
}

func (ns *NamespaceList) Init(whitelist []string, blacklist []string) {

	namespaceList, err := Client.GetNamespaces()
	if err != nil {
		panic(err)
	}

	ns.total = len(namespaceList.Items)
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
		ns.namespaces = append(ns.namespaces, *n)
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
