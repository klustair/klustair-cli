package trivyscanner

import (
	"math"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/klustair/cvssv3"
	log "github.com/sirupsen/logrus"
	cvssv2 "github.com/umisama/go-cvss"
)

type Cvss struct {
	V2 struct {
		Vector string `json:"vector"`
		Vendor string `json:"vendor"`
		Scores struct {
			Base          float64 `json:"base"`
			Temporal      float64 `json:"temporal"`
			Environmental float64 `json:"environmental"`
		} `json:"scores"`
		Metrics cvssv2.Vectors `json:"metrics"`
	} `json:"v2"`
	V3 struct {
		Vector string `json:"vector"`
		Vendor string `json:"vendor"`
		Scores struct {
			Base          float64 `json:"base"`
			Temporal      float64 `json:"temporal"`
			Environmental float64 `json:"environmental"`
			//Impact         float64 `json:"impact"`
			//Exploitability float64 `json:"exploitability"`
		} `json:"scores"`
		Metrics map[string]string `json:"metrics"`
	} `json:"v3"`

	Score   float64 `json:"score"`
	Version float64 `json:"version"`
}

/*
"CVSS": {
	"nvd": {
		"V2Vector": "AV:N/AC:L/Au:N/C:P/I:P/A:P",
		"V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
		"V2Score": 7.5,
		"V3Score": 9.8
	},
	"redhat": {
		"V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
		"V3Score": 9.8
	}
},
*/

/*
{
    "nvd": {
        "V2Score": 5,
        "V3Score": 7.5,
        "V2Vector": "AV:N/AC:L/Au:N/C:N/I:P/A:N",
        "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:H/A:N",
        "provider": "nvd",
        "V2Vector_metrics": {
            "A": "N",
            "C": "N",
            "I": "P",
            "AC": "L",
            "AV": "N",
            "Au": "N"
        },
        "V3Vector_metrics": {
            "A": "N",
            "C": "N",
            "I": "H",
            "S": "U",
            "AC": "L",
            "AV": "N",
            "MA": "N",
            "MC": "N",
            "MI": "H",
            "PR": "N",
            "UI": "N",
            "MAC": "L",
            "MAV": "N",
            "MPR": "N",
            "MUI": "N"
        },
        "V2Vector_base_score": "5.0",
        "V3Vector_base_score": "7.5",
        "V3Vector_modified_esc": "3.9",
        "V3Vector_modified_isc": "3.6"
    },
    "redhat": {
        "V2Score": 1.9,
        "V2Vector": "AV:L/AC:M/Au:N/C:N/I:P/A:N",
        "provider": "redhat",
        "V2Vector_metrics": {
            "A": "N",
            "C": "N",
            "I": "P",
            "AC": "M",
            "AV": "L",
            "Au": "N"
        },
        "V2Vector_base_score": "1.9"
    }
}
*/

func NewCVSS(CVSS types.VendorCVSS) *Cvss {

	for vendor, v := range CVSS {
		if v.V3Vector != "" && vendor == "nvd" {
			cvss := new(Cvss)
			return cvss.parseV3(v, vendor)
		}
		if v.V3Vector != "" && vendor == "redhat" {
			cvss := new(Cvss)
			return cvss.parseV3(v, vendor)
		} else if v.V2Vector != "" && vendor == "nvd" {
			cvss := new(Cvss)
			return cvss.parseV2(v, vendor)
		} else if v.V2Vector != "" && vendor == "redhat" {
			cvss := new(Cvss)
			return cvss.parseV2(v, vendor)
		}
		cvss := new(Cvss)
		cvss.empty("vendor:" + vendor)
		log.Debugf("Failed to parse CVSS vector from %s: %v", v, vendor)
		return cvss
	}
	log.Debugf("Failed to parse CVSS vector %v", CVSS)

	cvss := new(Cvss)
	cvss.empty("empty")
	return cvss
}

func (cvss *Cvss) empty(vendor string) *Cvss {
	cvss.V2.Vendor = vendor
	cvss.V2.Scores.Base = 0.0
	cvss.Version = 0.0
	cvss.Score = 0.0
	return cvss
}

func (cvss *Cvss) parseV2(v types.CVSS, vendor string) *Cvss {
	//log.Debugf("CVSS: %s %s %s", v.V2Vector, v.V2Score, vendor) //Too much Logoutput
	cvss.V2.Vector = v.V2Vector
	cvss.V2.Vendor = vendor

	v2Vector, err := cvssv2.ParseVectors("(" + v.V2Vector + ")")
	if err != nil {
		log.Infof("Failed to parse CVSSv2 vector: %s", err)
	}

	cvss.V2.Scores.Base = v2Vector.BaseScore()
	if v2Vector.HasTemporalVectors() && !math.IsNaN(v2Vector.TemporalScore()) {
		cvss.V2.Scores.Temporal = v2Vector.TemporalScore()
	}
	if v2Vector.HasEnvironmentalVectors() && !math.IsNaN(v2Vector.EnvironmentalScore()) {
		cvss.V2.Scores.Environmental = v2Vector.EnvironmentalScore()
	}
	cvss.V2.Metrics = v2Vector
	cvss.V2.Scores.Base = v.V2Score
	cvss.Version = 2.0
	return cvss
}

func (cvss *Cvss) parseV3(v types.CVSS, vendor string) *Cvss {
	//log.Debugf("CVSS: %s %s %s", v.V3Vector, v.V3Score, vendor) //Too much Logoutput
	cvss.V3.Vector = v.V3Vector
	cvss.V3.Vendor = vendor

	v3Vector, err := cvssv3.ParseVector(v.V3Vector)
	if err != nil {
		log.Infof("Failed to parse CVSSv3 vector: %s", err)
	}
	cvss.V3.Scores.Base = v3Vector.BaseScore()
	cvss.V3.Scores.Environmental = v3Vector.EnvironmentalScore()
	cvss.V3.Scores.Temporal = v3Vector.TemporalScore()
	cvss.V3.Metrics = v3Vector
	cvss.Version = 3.0
	cvss.Score = v.V3Score
	return cvss
}
