package klustair

import (
	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Pod struct {
	Podname            string    `json:"name"`
	Uid                string    `json:"uid"`
	Report_uid         string    `json:"report_uid"`
	Namespace_uid      string    `json:"namespace_uid"`
	Kubernetes_pod_uid types.UID `json:"kubernetes_pod_uid"`
	Creation_timestamp int64     `json:"creation_timestamp"`
	Age                int       `json:"age"`
}

func (p *Pod) Init(namespace_uid string, pod v1.Pod) {
	p.Podname = pod.Name
	p.Uid = uuid.New().String()
	p.Namespace_uid = namespace_uid
	p.Kubernetes_pod_uid = pod.UID
	p.Creation_timestamp = pod.CreationTimestamp.Unix()
}
