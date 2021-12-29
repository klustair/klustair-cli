package klustair

import (
	"fmt"
	"strings"

	"github.com/klustair/klustair-go/pkg/kubectl"
	"github.com/urfave/cli/v2"
)

type Options struct {
	Namespaces           []string
	NamespacesBlacklist  []string
	KubeAudit            []string
	Trivy                bool
	Label                string
	TrivyCredentialsPath string
	LimitDate            int
	LimitNr              int
	Configkey            string
	Apihost              string
	Apitoken             string
}

var Client *kubectl.Client

func RunCli(ctx *cli.Context) error {
	fmt.Println("run")
	opt, _ := loadOpts(ctx)
	Run(opt)
	return nil
	//return xerrors.Errorf("option error: %w", "nothing to do")
}

func Run(opt Options) error {
	fmt.Printf("run with options: %+v\n", opt)

	//initialize Kubectl client
	Client = kubectl.NewKubectlClient(false)

	//initialize Klustair Report
	Report := NewReport(opt)

	// TODO debug remove me
	for _, trivyreport := range Report.objectsList.trivyreports {
		fmt.Printf("trivyreport: %+v\n", trivyreport.ArtifactName)
	}

	// TODO debug remove me
	fmt.Printf("kubeauditReport: %+v\n", Report.kubeauditReports)
	fmt.Printf("Report: %+v\n", Report)

	if opt.Apihost != "" && opt.Apitoken != "" {
		Report.SendReport(opt)
	}
	return nil
	//return xerrors.Errorf("option error: %w", "nothing to do")
}

func loadOpts(ctx *cli.Context) (Options, error) {
	opt := Options{
		Namespaces:           strings.Split(ctx.String("namespaces"), ","),
		NamespacesBlacklist:  strings.Split(ctx.String("namespacesblacklist"), ","),
		KubeAudit:            strings.Split(ctx.String("kubeaudit"), ","),
		Trivy:                ctx.Bool("trivy"),
		Label:                ctx.String("label"),
		TrivyCredentialsPath: ctx.String("repocredentialspath"),
		LimitDate:            ctx.Int("limitdate"),
		LimitNr:              ctx.Int("limitnr"),
		Configkey:            ctx.String("configkey"),
		Apihost:              ctx.String("apihost"),
		Apitoken:             ctx.String("apitoken"),
	}
	return opt, nil
}
