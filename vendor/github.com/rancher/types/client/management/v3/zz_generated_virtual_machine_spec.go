package client

const (
	VirtualMachineSpecType             = "virtualMachineSpec"
	VirtualMachineSpecFieldDescription = "description"
	VirtualMachineSpecFieldProjectId   = "projectId"
)

type VirtualMachineSpec struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	ProjectId   string `json:"projectId,omitempty" yaml:"projectId,omitempty"`
}
