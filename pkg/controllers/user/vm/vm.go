package vm

import (
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
)

// TODO: Remove Warnf

type vmHandler struct {
}

func Register(cluster *config.UserContext) {
	logrus.Infof("Registering vm")

	vmH := &vmHandler{}
	//cluster.Management.Management.VirtualMachines("").AddClusterScopedHandler("vm", cluster.ClusterName, vmH.Sync)
	cluster.Management.Management.VirtualMachines("").AddClusterScopedLifecycle("vmHandler", cluster.ClusterName, vmH)
}

// Sync invokes the Policy Handler to take care of installing the native network policies
func (vmH *vmHandler) Sync(key string, vm *v3.VirtualMachine) error {
	if vm == nil {
		return nil
	}
	logrus.Warnf("vmH: updated vm=%+v", vm)
	return nil
}

func (vmH *vmHandler) Create(vm *v3.VirtualMachine) (*v3.VirtualMachine, error) {
	if vm == nil {
		return vm, nil
	}
	logrus.Warnf("vmHandler: Create: vm=%+v", *vm)
	return vm, nil
}

func (vmH *vmHandler) Updated(vm *v3.VirtualMachine) (*v3.VirtualMachine, error) {
	if vm == nil {
		return vm, nil
	}
	logrus.Warnf("vmHandler: Updated: vm=%+v", *vm)
	return vm, nil
}

func (vmH *vmHandler) Remove(vm *v3.VirtualMachine) (*v3.VirtualMachine, error) {
	if vm == nil {
		return vm, nil
	}
	logrus.Warnf("vmHandler: Remove: vm=%+v", *vm)
	return vm, nil
}
