/*
 * This file is part of the kiagnose project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
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
 * Copyright 2023 Red Hat, Inc.
 *
 */

package checkup

import (
	"github.com/kiagnose/kubevirt-storage-checkup/pkg/internal/checkup/vmi"

	corev1 "k8s.io/api/core/v1"

	kvcorev1 "kubevirt.io/api/core/v1"
)

const (
	guestMemory                   = "2Gi"
	rootDiskName                  = "rootdisk"
	terminationGracePeriodSeconds = 0
)

func newVMUnderTest(name string, pvc *corev1.PersistentVolumeClaim) *kvcorev1.VirtualMachine {
	optionsToApply := []vmi.Option{
		vmi.WithDataVolume(rootDiskName, pvc),
		vmi.WithMemory(guestMemory),
		vmi.WithTerminationGracePeriodSeconds(terminationGracePeriodSeconds),
	}
	return vmi.NewVM(name, optionsToApply...)
}
