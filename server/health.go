// Copyright 2022 Leon Ding <ding@ibyte.me> https://wiredb.github.io

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type health struct {
	mem  *mem.VirtualMemoryStat
	disk *disk.UsageStat
}

func newHealth(path string) (*health, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	disk, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	return &health{mem: mem, disk: disk}, nil
}

// GetTotalMemory returns the total system memory in bytes.
func (h *health) GetTotalMemory() uint64 {
	return h.mem.Total
}

// GetFreeMemory returns the available system memory in bytes.
func (h *health) GetFreeMemory() uint64 {
	return h.mem.Available
}

func (h *health) GetUsedDisk() uint64 {
	return h.disk.Used
}

func (h *health) GetFreeDisk() uint64 {
	return h.disk.Free
}

func (h *health) GetTotalDisk() uint64 {
	return h.disk.Total
}

func (h *health) GetDiskPercent() float64 {
	return h.disk.UsedPercent
}
