package webhooks

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfield "k8s.io/apimachinery/pkg/util/validation/field"

	v1 "kubevirt.io/api/core/v1"
)

var _false bool = false

const (
	defaultCPUModelRiscv64 = v1.CPUModeHostPassthrough
)
// 仿照arm64.go写的，不知道是不是真的需要验证这些项目
// ValidateVirtualMachineInstanceRiscv64Setting is a validation function for validating-webhook to filter unsupported setting on Riscv64
func ValidateVirtualMachineInstanceRiscv64Setting(field *k8sfield.Path, spec *v1.VirtualMachineInstanceSpec) []metav1.StatusCause {
	var statusCauses []metav1.StatusCause
	validateBootOptions(field, spec, &statusCauses)
	validateCPUModel(field, spec, &statusCauses)
	validateDiskBus(field, spec, &statusCauses)
	validateWatchdog(field, spec, &statusCauses)
	validateSoundDevice(field, spec, &statusCauses)
	return statusCauses
}

func validateBootOptions(field *k8sfield.Path, spec *v1.VirtualMachineInstanceSpec, statusCauses *[]metav1.StatusCause) {
	if spec.Domain.Firmware != nil && spec.Domain.Firmware.Bootloader != nil {
		if spec.Domain.Firmware.Bootloader.BIOS != nil {
			*statusCauses = append(*statusCauses, metav1.StatusCause{
				Type:    metav1.CauseTypeFieldValueNotSupported,
				Message: "Riscv64 does not support BIOS boot, please change to UEFI boot",
				Field:   field.Child("domain", "firmware", "bootloader", "bios").String(),
			})
		}
		if spec.Domain.Firmware.Bootloader.EFI != nil {
			if spec.Domain.Firmware.Bootloader.EFI.SecureBoot == nil || (spec.Domain.Firmware.Bootloader.EFI.SecureBoot != nil && *spec.Domain.Firmware.Bootloader.EFI.SecureBoot) {
				*statusCauses = append(*statusCauses, metav1.StatusCause{
					Type:    metav1.CauseTypeFieldValueNotSupported,
					Message: "UEFI secure boot is currently not supported on riscv64 architecture",
					Field:   field.Child("domain", "firmware", "bootloader", "efi", "secureboot").String(),
				})
			}
		}
	}
}

func validateCPUModel(field *k8sfield.Path, spec *v1.VirtualMachineInstanceSpec, statusCauses *[]metav1.StatusCause) {
	if spec.Domain.CPU != nil && (&spec.Domain.CPU.Model != nil) && spec.Domain.CPU.Model == "host-model" {
		*statusCauses = append(*statusCauses, metav1.StatusCause{
			Type:    metav1.CauseTypeFieldValueNotSupported,
			Message: "Riscv64 does not support CPU host-model",
			Field:   field.Child("domain", "cpu", "model").String(),
		})
	}
}

func validateDiskBus(field *k8sfield.Path, spec *v1.VirtualMachineInstanceSpec, statusCauses *[]metav1.StatusCause) {
	if spec.Domain.Devices.Disks != nil {
		// checkIfBusAvailable: if bus type is nil, virtio, scsi return true, otherwise, return false
		checkIfBusAvailable := func(bus v1.DiskBus) bool {
			if bus == "" || bus == v1.DiskBusVirtio || bus == v1.DiskBusSCSI {
				return true
			}
			return false
		}

		for i, disk := range spec.Domain.Devices.Disks {
			if disk.Disk != nil && !checkIfBusAvailable(disk.Disk.Bus) {
				*statusCauses = append(*statusCauses, metav1.StatusCause{
					Type:    metav1.CauseTypeFieldValueNotSupported,
					Message: "Riscv64 does not support this disk bus type, please use virtio or scsi",
					Field:   field.Child("domain", "devices", "disks").Index(i).Child("disk", "bus").String(),
				})
			}
			if disk.CDRom != nil && !checkIfBusAvailable(disk.CDRom.Bus) {
				*statusCauses = append(*statusCauses, metav1.StatusCause{
					Type:    metav1.CauseTypeFieldValueNotSupported,
					Message: "Riscv64 does not support this disk bus type, please use virtio or scsi",
					Field:   field.Child("domain", "devices", "disks").Index(i).Child("cdrom", "bus").String(),
				})
			}
			if disk.LUN != nil && !checkIfBusAvailable(disk.LUN.Bus) {
				*statusCauses = append(*statusCauses, metav1.StatusCause{
					Type:    metav1.CauseTypeFieldValueNotSupported,
					Message: "Riscv64 does not support this disk bus type, please use virtio or scsi",
					Field:   field.Child("domain", "devices", "disks").Index(i).Child("lun", "bus").String(),
				})
			}
		}
	}
}

func validateWatchdog(field *k8sfield.Path, spec *v1.VirtualMachineInstanceSpec, statusCauses *[]metav1.StatusCause) {
	if spec.Domain.Devices.Watchdog != nil {
		*statusCauses = append(*statusCauses, metav1.StatusCause{
			Type:    metav1.CauseTypeFieldValueNotSupported,
			Message: "Riscv64 does not support Watchdog device",
			Field:   field.Child("domain", "devices", "watchdog").String(),
		})
	}
}

func validateSoundDevice(field *k8sfield.Path, spec *v1.VirtualMachineInstanceSpec, statusCauses *[]metav1.StatusCause) {
	if spec.Domain.Devices.Sound != nil {
		*statusCauses = append(*statusCauses, metav1.StatusCause{
			Type:    metav1.CauseTypeFieldValueNotSupported,
			Message: "Riscv64 does not support sound device",
			Field:   field.Child("domain", "devices", "sound").String(),
		})
	}
}

