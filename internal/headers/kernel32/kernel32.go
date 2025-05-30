// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows

package kernel32

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

//nolint:gochecknoglobals
var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetDynamicTimeZoneInformationSys = modkernel32.NewProc("GetDynamicTimeZoneInformation")
	procKernelLocalFileTimeToFileTime    = modkernel32.NewProc("LocalFileTimeToFileTime")
	procGetTickCount                     = modkernel32.NewProc("GetTickCount64")
)

// SYSTEMTIME contains a date and time.
// 📑 https://docs.microsoft.com/en-us/windows/win32/api/minwinbase/ns-minwinbase-systemtime
type SYSTEMTIME struct {
	WYear         uint16
	WMonth        uint16
	WDayOfWeek    uint16
	WDay          uint16
	WHour         uint16
	WMinute       uint16
	WSecond       uint16
	WMilliseconds uint16
}

// DynamicTimezoneInformation contains the current dynamic daylight time settings.
// 📑 https://docs.microsoft.com/en-us/windows/win32/api/timezoneapi/ns-timezoneapi-dynamic_time_zone_information
type DynamicTimezoneInformation struct {
	Bias                        int32
	standardName                [32]uint16
	StandardDate                SYSTEMTIME
	StandardBias                int32
	DaylightName                [32]uint16
	DaylightDate                SYSTEMTIME
	DaylightBias                int32
	TimeZoneKeyName             [128]uint16
	DynamicDaylightTimeDisabled uint8 // BOOLEAN
}

// GetDynamicTimeZoneInformation retrieves the current dynamic daylight time settings.
// 📑 https://docs.microsoft.com/en-us/windows/win32/api/timezoneapi/nf-timezoneapi-getdynamictimezoneinformation
func GetDynamicTimeZoneInformation() (DynamicTimezoneInformation, error) {
	var tzi DynamicTimezoneInformation

	r0, _, err := procGetDynamicTimeZoneInformationSys.Call(uintptr(unsafe.Pointer(&tzi)))
	if uint32(r0) == 0xffffffff {
		return tzi, err
	}

	return tzi, nil
}

func LocalFileTimeToFileTime(localFileTime, utcFileTime *windows.Filetime) uint32 {
	ret, _, _ := procKernelLocalFileTimeToFileTime.Call(
		uintptr(unsafe.Pointer(localFileTime)),
		uintptr(unsafe.Pointer(utcFileTime)))

	return uint32(ret)
}

func GetTickCount64() uint64 {
	ret, _, _ := procGetTickCount.Call()

	return uint64(ret)
}
