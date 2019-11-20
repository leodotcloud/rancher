package client

const (
	ClusterScanConfigType               = "clusterScanConfig"
	ClusterScanConfigFieldCISScanConfig = "cisScanConfig"
)

type ClusterScanConfig struct {
	CISScanConfig *CisScanConfig `json:"cisScanConfig,omitempty" yaml:"cisScanConfig,omitempty"`
}
