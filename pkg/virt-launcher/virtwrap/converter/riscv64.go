package converter

import (
	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/kubevirt/pkg/pointer"
	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
)

// Ensure that there is a compile error should the struct not implement the archConverter interface anymore.
var _ = archConverter(&archConverterRISCV64{})

type archConverterRISCV64 struct{}

func (archConverterRISCV64) addGraphicsDevice(vmi *v1.VirtualMachineInstance, domain *api.Domain, _ *ConverterContext) {
	// For riscv64, qemu-kvm only supports virtio-gpu display device, so set it as default video device.
	// tablet and keyboard devices are necessary for controlling the VM via vnc connection
	domain.Spec.Devices.Video = []api.Video{
		{
			Model: api.VideoModel{
				Type:  v1.VirtIO,
				Heads: pointer.P(graphicsDeviceDefaultHeads),
			},
		},
	}

	if !hasTabletDevice(vmi) {
		domain.Spec.Devices.Inputs = append(domain.Spec.Devices.Inputs,
			api.Input{
				Bus:  "usb",
				Type: "tablet",
			},
		)
	}

	domain.Spec.Devices.Inputs = append(domain.Spec.Devices.Inputs,
		api.Input{
			Bus:  "usb",
			Type: "keyboard",
		},
	)
}

func (archConverterRISCV64) scsiController(c *ConverterContext, driver *api.ControllerDriver) api.Controller {
	return defaultSCSIController(c, driver)
}

func (archConverterRISCV64) isUSBNeeded(_ *v1.VirtualMachineInstance) bool {
	return true
}

