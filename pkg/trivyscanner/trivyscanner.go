package trivyscanner

import (
	"os"
	"time"

	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/commands/option"
	"github.com/klustair/trivy/pkg/commands/artifact"
)

func GetDefaultOption() (artifact.Option, error) {
	opt := artifact.Option{
		GlobalOption: option.GlobalOption{
			Context:    nil,
			Logger:     nil,
			AppVersion: "1",
			Quiet:      false,
			Debug:      true,
			CacheDir:   "/tmp/trivy",
		},
		ArtifactOption: option.ArtifactOption{
			Input:      "",
			Timeout:    time.Duration(15 * time.Second),
			ClearCache: false,

			SkipDirs:  []string{},
			SkipFiles: []string{},

			Target: "node:latest",
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
			Format:   "table",
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
				//"config",
			},
			Output: os.Stdout, //nil, //os.Stdout, //file
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

	return opt, nil
}
