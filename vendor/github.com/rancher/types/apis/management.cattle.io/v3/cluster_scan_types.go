package v3

import (
	"github.com/rancher/norman/condition"
	"github.com/rancher/norman/types"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ClusterScanConditionCreated   condition.Cond = "Created"
	ClusterScanConditionCompleted condition.Cond = "Completed"

	ClusterScanTypeCIS     = "cis"
	DefaultNamespaceForCis = "security-scan"
	DefaultSonobuoyPodName = "security-scan-runner"

	RunCISScanAnnotation         = "field.cattle.io/runCisScan"
	SonobuoyCompletionAnnotation = "field.cattle.io/sonobuoyDone"
	CisHelmChartOwner            = "field.cattle.io/clusterScanOwner"
)

type CisScanConfig struct {
	// Show only checks which include failures in the final report
	FailuresOnly bool `json:"failuresOnly"`
	// IDs of the checks that need to be skipped in the final report
	Skip []string `json:"skip"`
}

type ClusterScanConfig struct {
	// Config related to CIS Scan type
	CISScanConfig *CisScanConfig `json:"cisScanConfig"`
}

type ClusterScanCondition struct {
	// Type of condition.
	Type string `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition
	Message string `json:"message,omitempty"`
}

type ClusterScanSpec struct {
	ScanType string `json:"scanType"`
	// cluster ID
	ClusterID string `json:"clusterId,omitempty" norman:"required,type=reference[cluster]"`
	// manual flag
	Manual bool `yaml:"manual" json:"manual,omitempty"`
	// scanConfig
	ScanConfig ClusterScanConfig `yaml:",omitempty" json:"scanConfig,omitempty"`
}

type ClusterScanStatus struct {
	Conditions []ClusterScanCondition `json:"conditions"`
}

type ClusterScan struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterScanSpec   `json:"spec"`
	Status ClusterScanStatus `yaml:"status" json:"status,omitempty"`
}
