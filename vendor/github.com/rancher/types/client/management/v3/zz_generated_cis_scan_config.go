package client

const (
	CisScanConfigType              = "cisScanConfig"
	CisScanConfigFieldFailuresOnly = "failuresOnly"
	CisScanConfigFieldSkip         = "skip"
)

type CisScanConfig struct {
	FailuresOnly bool     `json:"failuresOnly,omitempty" yaml:"failuresOnly,omitempty"`
	Skip         []string `json:"skip,omitempty" yaml:"skip,omitempty"`
}
