package kubeaudit

import (
	"time"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	kubeauditconfig "github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

//var KubeauditReport *kubeaudit.Report

//type KubeauditReport *kubeaudit.Report

type Auditor struct {
	KubeauditConfig kubeauditconfig.KubeauditConfig
	Report          kubeaudit.Report
	Klustair        struct {
		ReportUid    string `json:"report_uid"`
		NamespaceUid string `json:"namespace_uid"`
	}
}

func (a *Auditor) SetConfig(auditors []string) kubeauditconfig.KubeauditConfig {
	auditoorsmap := make(map[string]bool)
	for _, a := range auditors {
		log.Debugf("auditor: %+v\n", a)
		auditoorsmap[a] = true
	}
	a.KubeauditConfig.EnabledAuditors = auditoorsmap

	return a.KubeauditConfig
}

func (a *Auditor) Audit(namespace string) []KubeauditReport {
	auditors, err := all.Auditors(kubeauditconfig.KubeauditConfig{})
	if err != nil {
		panic(err)
	}

	kubeAuditor, err := kubeaudit.New(auditors)
	if err != nil {
		panic(err)
	}

	// TODO Need some love here.
	if true {
		report, err := kubeAuditor.AuditLocal("", kubeaudit.AuditOptions{Namespace: namespace})
		if err != nil {
			panic(err)
		}
		a.Report = *report

		return a.getReport()
	} else {
		report, err := kubeAuditor.AuditCluster(kubeaudit.AuditOptions{Namespace: namespace})
		if err != nil {
			panic(err)
		}
		a.Report = *report

		return a.getReport()
	}
}

type KubeauditReport struct {
	Uid                string `json:"uid"`
	ReportUid          string `json:"report_uid"`
	NamespaceUid       string `json:"namespace_uid"`
	AuditTime          string `json:"time"`
	AuditType          string `json:"audit_type"`
	AuditName          string `json:"AuditResultName"`
	Message            string `json:"msg"`
	SeverityLevel      string `json:"level"`
	Capability         string `json:"Capability"`
	Container          string `json:"Container"`
	MissingAnnotations string `json:"MissingAnnotations"`
	ResourceName       string `json:"ResourceName"`
	ResourceNamespace  string `json:"ResourceNamespace"`
	ResourceApiVersion string `json:"ResourceApiVersion"`
}

func (a *Auditor) getReport() []KubeauditReport {
	//a.Report.PrintResults()
	var kubeauditReports []KubeauditReport
	for _, r := range a.Report.Results() {
		//log.Debugf("Result: %+v\n", r)
		ok := r.GetResource().Object().GetObjectKind()
		//objectGroup := ok.GroupVersionKind().Group
		objectKind := ok.GroupVersionKind().Kind
		objectMeta := k8s.GetObjectMeta(r.GetResource().Object())

		resourceApiVersion := ok.GroupVersionKind().Version
		resourceName := objectMeta.GetName()
		resourceNamespace := objectMeta.GetNamespace()

		for _, o := range r.GetAuditResults() {

			k := KubeauditReport{
				Uid:                uuid.New().String(),
				ReportUid:          a.Klustair.ReportUid,
				NamespaceUid:       a.Klustair.NamespaceUid,
				AuditName:          o.Name,
				Message:            o.Message,
				AuditTime:          time.Now().UTC().Format(time.RFC3339),
				SeverityLevel:      o.Severity.String(),
				AuditType:          "unknown", //initial value
				ResourceName:       resourceName,
				ResourceNamespace:  resourceNamespace,
				ResourceApiVersion: resourceApiVersion,
			}

			if o.Name == "CapabilityAddedAllowed" {
				k.Capability = o.Metadata["Metadata"]
			}

			// ugly but need to flatten this out a bit
			if objectKind == "Deployment" ||
				objectKind == "StatefulSet" ||
				objectKind == "DaemonSet" ||
				objectKind == "ReplicaSet" ||
				objectKind == "Pod" ||
				objectKind == "ReplicationController" ||
				objectKind == "Job" ||
				objectKind == "PodTemplate" ||
				objectKind == "CronJob" {
				k.AuditType = "pod"
			} else if objectKind == "Namespace" {
				k.AuditType = "namespace"
			} else if objectKind == "Service" {
				k.AuditType = "service"
			} else if objectKind == "ServiceAccount" {
				k.AuditType = "serviceAccount"
			} else if objectKind == "NetworkPolicy" {
				k.AuditType = "networkPolicy"
			}

			if container, ok := o.Metadata["Container"]; ok {
				k.AuditType = "container"
				k.Container = container
			}

			if annotation, ok := o.Metadata["MissingAnnotation"]; ok {
				k.MissingAnnotations = annotation
			}
			//log.Debugf("KubeauditReport: %+v\n", k)
			kubeauditReports = append(kubeauditReports, k)
		}
	}

	return kubeauditReports
}
