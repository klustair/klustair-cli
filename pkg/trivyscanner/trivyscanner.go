package trivyscanner

import (
	"context"
	"fmt"
	"os"
	"time"

	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/commands/option"
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/klustair/trivy/pkg/commands/artifact"
)

type Trivy struct {
	inCluster bool
	Options   artifact.Option
}

func (t *Trivy) NewScanner() *Trivy {
	return &Trivy{
		inCluster: false,
		Options:   GetOption(),
	}
}

func GetOption() artifact.Option {

	//file, _ := os.Create("/tmp/trivy.txt")
	//file, _ := os.Create(os.DevNull)

	trivy_debug, exists := os.LookupEnv("TRIVY_DEBUG")
	var debug bool
	debug = false
	if exists {
		if trivy_debug == "true" {
			debug = true
		}
	}
	fmt.Println("TRIVY_DEBUG:", debug)

	trivy_quiet, exists := os.LookupEnv("TRIVY_QUIET")
	var quiet bool
	quiet = true
	if exists {
		if trivy_quiet == "false" {
			quiet = false
		}
	}
	fmt.Println("TRIVY_QUIET:", quiet)

	option := artifact.Option{
		GlobalOption: option.GlobalOption{
			Context:    nil,
			Logger:     nil,
			AppVersion: "1",
			Quiet:      quiet,
			Debug:      debug,
			CacheDir:   "/tmp/trivy",
		},
		ArtifactOption: option.ArtifactOption{
			Input:      "",
			Timeout:    time.Duration(150 * time.Second),
			ClearCache: false,

			SkipDirs:  []string{},
			SkipFiles: []string{},

			Target:      "node:latest",
			OfflineScan: false,
		},
		DBOption: option.DBOption{
			Reset:          false,
			DownloadDBOnly: false,
			SkipDBUpdate:   false,
			Light:          false,
			NoProgress:     false,
		},
		ImageOption: option.ImageOption{
			ScanRemovedPkgs: false,
			ListAllPkgs:     false,
		},
		ReportOption: option.ReportOption{
			Format:   "template",
			Template: "",

			IgnoreFile:    "",
			IgnoreUnfixed: false,
			ExitCode:      0,
			IgnorePolicy:  "",

			VulnType: []string{
				"os",
				"library",
			},
			SecurityChecks: []string{
				"vuln",
				"config",
			},
			Output: nil, //nil, //os.Stdout, //file
			Severities: []dbTypes.Severity{
				0,
				1,
				2,
				3,
				4,
			},
		},
		CacheOption: option.CacheOption{
			CacheBackend: "fs",
		},
		ConfigOption: option.ConfigOption{
			FilePatterns:       nil,
			IncludeNonFailures: false,
			SkipPolicyUpdate:   true,
			Trace:              false,

			PolicyPaths:      []string{},
			DataPaths:        []string{},
			PolicyNamespaces: []string{},
		},
		DisabledAnalyzers: nil,
	}
	return option
}

func (t *Trivy) Scan(image string) (report.Report, error) {
	t.Options.ArtifactOption.Target = image

	// TODO: make timeout configurable
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	return artifact.ImageRunLib(ctx, t.Options)

}
