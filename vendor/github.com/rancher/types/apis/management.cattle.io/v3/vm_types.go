package v3

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualMachineSpec struct {
	ProjectName string `json:"projectName,omitempty" norman:"required,type=reference[project]"`
	Description string `json:"description"`
}

type VirtualMachineStatus struct {
}

type VirtualMachine struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VirtualMachineSpec    `json:"spec"`
	Status            *VirtualMachineStatus `json:"status"`
}
