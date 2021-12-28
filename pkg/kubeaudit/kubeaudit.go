package kubeaudit

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	kubeauditconfig "github.com/Shopify/kubeaudit/config"
)

//var KubeauditReport *kubeaudit.Report

type KubeauditReport *kubeaudit.Report

type Auditor struct {
	KubeauditConfig kubeauditconfig.KubeauditConfig
	Report          *kubeaudit.Report
}

func (a *Auditor) SetConfig(auditors string) kubeauditconfig.KubeauditConfig {
	auditoorsmap := make(map[string]bool)
	fmt.Printf("auditors: %+v\n", auditors)
	for _, a := range strings.Split(auditors, ",") {
		fmt.Printf("auditor: %+v\n", a)
		auditoorsmap[a] = true
	}
	a.KubeauditConfig.EnabledAuditors = auditoorsmap

	return a.KubeauditConfig
}

func (a *Auditor) Run() *kubeaudit.Report {
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
		report, err := kubeAuditor.AuditLocal("", kubeaudit.AuditOptions{})

		if err != nil {
			panic(err)
		}

		a.Report = report

		return report
	} else {
		report, err := kubeAuditor.AuditCluster(kubeaudit.AuditOptions{})

		if err != nil {
			panic(err)
		}

		a.Report = report

		return report
	}
}
