package client

import (
	"github.com/rancher/norman/types"
)

const (
	VirtualMachineType                      = "virtualMachine"
	VirtualMachineFieldAnnotations          = "annotations"
	VirtualMachineFieldCreated              = "created"
	VirtualMachineFieldCreatorID            = "creatorId"
	VirtualMachineFieldDescription          = "description"
	VirtualMachineFieldLabels               = "labels"
	VirtualMachineFieldName                 = "name"
	VirtualMachineFieldNamespaceId          = "namespaceId"
	VirtualMachineFieldOwnerReferences      = "ownerReferences"
	VirtualMachineFieldProjectId            = "projectId"
	VirtualMachineFieldRemoved              = "removed"
	VirtualMachineFieldState                = "state"
	VirtualMachineFieldStatus               = "status"
	VirtualMachineFieldTransitioning        = "transitioning"
	VirtualMachineFieldTransitioningMessage = "transitioningMessage"
	VirtualMachineFieldUuid                 = "uuid"
)

type VirtualMachine struct {
	types.Resource
	Annotations          map[string]string     `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string                `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string                `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description          string                `json:"description,omitempty" yaml:"description,omitempty"`
	Labels               map[string]string     `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string                `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId          string                `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference      `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ProjectId            string                `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Removed              string                `json:"removed,omitempty" yaml:"removed,omitempty"`
	State                string                `json:"state,omitempty" yaml:"state,omitempty"`
	Status               *VirtualMachineStatus `json:"status,omitempty" yaml:"status,omitempty"`
	Transitioning        string                `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string                `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	Uuid                 string                `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type VirtualMachineCollection struct {
	types.Collection
	Data   []VirtualMachine `json:"data,omitempty"`
	client *VirtualMachineClient
}

type VirtualMachineClient struct {
	apiClient *Client
}

type VirtualMachineOperations interface {
	List(opts *types.ListOpts) (*VirtualMachineCollection, error)
	Create(opts *VirtualMachine) (*VirtualMachine, error)
	Update(existing *VirtualMachine, updates interface{}) (*VirtualMachine, error)
	ByID(id string) (*VirtualMachine, error)
	Delete(container *VirtualMachine) error
}

func newVirtualMachineClient(apiClient *Client) *VirtualMachineClient {
	return &VirtualMachineClient{
		apiClient: apiClient,
	}
}

func (c *VirtualMachineClient) Create(container *VirtualMachine) (*VirtualMachine, error) {
	resp := &VirtualMachine{}
	err := c.apiClient.Ops.DoCreate(VirtualMachineType, container, resp)
	return resp, err
}

func (c *VirtualMachineClient) Update(existing *VirtualMachine, updates interface{}) (*VirtualMachine, error) {
	resp := &VirtualMachine{}
	err := c.apiClient.Ops.DoUpdate(VirtualMachineType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *VirtualMachineClient) List(opts *types.ListOpts) (*VirtualMachineCollection, error) {
	resp := &VirtualMachineCollection{}
	err := c.apiClient.Ops.DoList(VirtualMachineType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *VirtualMachineCollection) Next() (*VirtualMachineCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &VirtualMachineCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *VirtualMachineClient) ByID(id string) (*VirtualMachine, error) {
	resp := &VirtualMachine{}
	err := c.apiClient.Ops.DoByID(VirtualMachineType, id, resp)
	return resp, err
}

func (c *VirtualMachineClient) Delete(container *VirtualMachine) error {
	return c.apiClient.Ops.DoResourceDelete(VirtualMachineType, &container.Resource)
}
