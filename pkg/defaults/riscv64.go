/* Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2021
 *
 */
package defaults

import (
	v1 "kubevirt.io/api/core/v1"
)

var _false_rv bool = false

const (
	defaultCPUModelRISCVC64 = v1.CPUModeHostPassthrough
)

// setDefaultRISCVC64CPUModel set default CPU model to host-passthrough for RISCV
func setDefaultRISCVC64CPUModel(spec *v1.VirtualMachineInstanceSpec) {
	if spec.Domain.CPU == nil {
		spec.Domain.CPU = &v1.CPU{}
	}

	if spec.Domain.CPU.Model == "" {
		spec.Domain.CPU.Model = defaultCPUModelRISCVC64
	}
}

// setDefaultRISCVC64Bootloader set default bootloader to UEFI boot for RISCV
func setDefaultRISCVC64Bootloader(spec *v1.VirtualMachineInstanceSpec) {
	if spec.Domain.Firmware == nil || spec.Domain.Firmware.Bootloader == nil {
		if spec.Domain.Firmware == nil {
			spec.Domain.Firmware = &v1.Firmware{}
		}
		if spec.Domain.Firmware.Bootloader == nil {
			spec.Domain.Firmware.Bootloader = &v1.Bootloader{}
		}
		spec.Domain.Firmware.Bootloader.EFI = &v1.EFI{}
		spec.Domain.Firmware.Bootloader.EFI.SecureBoot = &_false_rv
	}
}

// setDefaultRISCVC64DisksBus set default Disks Bus for RISCV, as QEMU-KVM for RISCV might not support SATA
func setDefaultRISCVC64DisksBus(spec *v1.VirtualMachineInstanceSpec) {
	bus := v1.DiskBusVirtio

	for i := range spec.Domain.Devices.Disks {
		disk := &spec.Domain.Devices.Disks[i].DiskDevice

		if disk.Disk != nil && disk.Disk.Bus == "" {
			disk.Disk.Bus = bus
		}
		if disk.CDRom != nil && disk.CDRom.Bus == "" {
			disk.CDRom.Bus = bus
		}
		if disk.LUN != nil && disk.LUN.Bus == "" {
			disk.LUN.Bus = bus
		}
	}
}

// SetRISCVC64Defaults is mutating function for mutating-webhook for RISCV
func SetRISCVC64Defaults(spec *v1.VirtualMachineInstanceSpec) {
	setDefaultRISCVC64CPUModel(spec)
	setDefaultRISCVC64Bootloader(spec)
	setDefaultRISCVC64DisksBus(spec)
}

// IsRISCVC64 checks if the architecture is RISCV
func IsRISCVC64(vmiSpec *v1.VirtualMachineInstanceSpec) bool {
	return vmiSpec.Architecture == "riscv64"
}

