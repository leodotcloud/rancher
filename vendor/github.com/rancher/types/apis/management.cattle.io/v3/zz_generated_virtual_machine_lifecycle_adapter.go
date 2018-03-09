package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type VirtualMachineLifecycle interface {
	Create(obj *VirtualMachine) (*VirtualMachine, error)
	Remove(obj *VirtualMachine) (*VirtualMachine, error)
	Updated(obj *VirtualMachine) (*VirtualMachine, error)
}

type virtualMachineLifecycleAdapter struct {
	lifecycle VirtualMachineLifecycle
}

func (w *virtualMachineLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*VirtualMachine))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *virtualMachineLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*VirtualMachine))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *virtualMachineLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*VirtualMachine))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewVirtualMachineLifecycleAdapter(name string, clusterScoped bool, client VirtualMachineInterface, l VirtualMachineLifecycle) VirtualMachineHandlerFunc {
	adapter := &virtualMachineLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *VirtualMachine) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
