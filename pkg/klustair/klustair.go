package klustair

import (
	"context"
	"fmt"

	"github.com/klustair/klustair-go/pkg/kubectl"
	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	configkey            string
}

func RunCli(ctx *cli.Context) error {
	fmt.Println("run")
	opt, _ := loadOpts(ctx)
	Run(opt)
	return nil
	//return xerrors.Errorf("option error: %w", "nothing to do")
}

func Run(opt Options) error {
	fmt.Printf("run with options: %+v\n", opt)

	clientset, _ := kubectl.GetClientset(false)
	fmt.Printf("clientset: %+v\n", clientset)

	pods, err := clientset.CoreV1().Pods("klustair").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	return nil
	//return xerrors.Errorf("option error: %w", "nothing to do")
}

func loadOpts(ctx *cli.Context) (Options, error) {
	opt := Options{
		Namespaces:           ctx.StringSlice("namespaces"),
		NamespacesBlacklist:  ctx.StringSlice("namespacesblacklist"),
		KubeAudit:            ctx.StringSlice("kubeaudit"),
		Trivy:                ctx.Bool("trivy"),
		Label:                ctx.String("label"),
		TrivyCredentialsPath: ctx.String("repocredentialspath"),
		LimitDate:            ctx.Int("limitdate"),
		LimitNr:              ctx.Int("limitnr"),
		configkey:            ctx.String("configkey"),
	}
	return opt, nil
}
