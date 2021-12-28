package klustair

import (
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

func (p *Pod) Init(namespace_uid string, pod v1.Pod) {
	p.podname = pod.Name
	p.uid = uuid.New().String()
	p.namespace_uid = namespace_uid
	p.kubernetes_pod_uid = pod.UID
	p.creation_timestamp = pod.CreationTimestamp.Unix()
}
