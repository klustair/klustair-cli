package kubeaudit

import (
	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	kubeauditconfig "github.com/Shopify/kubeaudit/config"
)

func All() {
	auditors, err := all.Auditors(kubeauditconfig.KubeauditConfig{})
	if err != nil {
		panic(err)
	}

	kubeAuditor, err := kubeaudit.New(auditors)
	if err != nil {
		panic(err)
	}

	report, err := kubeAuditor.AuditCluster(kubeaudit.AuditOptions{})
	if err != nil {
		panic(err)
	}
	report.PrintResults()
}
