package klustair

import (
	"time"

	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Pod struct {
	Uid                string    `json:"uid"`
	ReportUid          string    `json:"report_uid"`
	NamespaceUid       string    `json:"namespace_uid"`
	Podname            string    `json:"podname"`
	Kubernetes_pod_uid types.UID `json:"kubernetes_pod_uid"`
	Creation_timestamp string    `json:"creation_timestamp"`
	Age                int       `json:"age"`
}

func (p *Pod) Init(reportUid string, namespaceUid string, pod v1.Pod) {
	p.Uid = uuid.New().String()
	p.ReportUid = reportUid
	p.NamespaceUid = namespaceUid
	p.Podname = pod.Name
	p.Kubernetes_pod_uid = pod.UID
	p.Creation_timestamp = pod.CreationTimestamp.UTC().Format(time.RFC3339)

	p.Age = int(time.Since(time.Unix(pod.CreationTimestamp.Unix(), 0)).Hours() / 24)
}
