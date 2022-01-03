package kubeaudit

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	kubeauditconfig "github.com/Shopify/kubeaudit/config"
	log "github.com/sirupsen/logrus"
)

//var KubeauditReport *kubeaudit.Report

//type KubeauditReport *kubeaudit.Report

type Auditor struct {
	KubeauditConfig kubeauditconfig.KubeauditConfig
	Report          KubeauditReport
}

type KubeauditReport struct {
	Uid                string `json:"uid"`
	ReportUid          string `json:"report_uid"`
	NamespaceUid       string `json:"namespace_uid"`
	AuditType          string `json:"audit_type"`
	AuditName          string `json:"AuditName"`
	Message            string `json:"msg"`
	SeverityLevel      string `json:"level"`
	ResourceName       string `json:"ResourceName"`
	Capability         string `json:"Capability"`
	Container          string `json:"Container"`
	AuditResultName    string `json:"AuditResultName"`
	MissingAnnotations string `json:"MissingAnnotations"`
	ResourceNamespace  string `json:"ResourceNamespace"`
	ResourceApiVersion string `json:"ResourceApiVersion"`
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

func (a *Auditor) Run(namespaces []string) []*kubeaudit.Report {
	var reports []*kubeaudit.Report
	for _, namespace := range namespaces {
		log.Debugf("Kubeaudit on namespace: %+v", namespace)
		report := a.Audit(namespace)
		reports = append(reports, report)
	}
	return reports
}

func (a *Auditor) Audit(namespace string) *kubeaudit.Report {
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

		a.Report = *a.getReport()

		return report
	} else {
		report, err := kubeAuditor.AuditCluster(kubeaudit.AuditOptions{Namespace: namespace})

		if err != nil {
			panic(err)
		}

		a.Report = *a.getReport()

		return report
	}
}

func (a *Auditor) getReport() *KubeauditReport {
	newReport := new(KubeauditReport)
	// TODO parse report
	newReport.Uid = "test"
	return newReport
}
