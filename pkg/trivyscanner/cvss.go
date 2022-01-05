package trivyscanner

type Cvss struct {
	BaseScore    float64 `json:"base_score"`
	BaseSeverity string  `json:"base_severity"`
	Vector       string  `json:"vector"`
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
