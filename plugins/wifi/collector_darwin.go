//go:build darwin && cgo

package wifi

/*
#cgo LDFLAGS: -framework CoreWLAN -framework CoreLocation -framework Foundation
#include <stdlib.h>
char *cmdry_wifi_info(void);
*/
import "C"

import (
	"strings"
	"unsafe"
)

// Collect reads the active connection through CoreWLAN rather than shelling out
// to networksetup, which may require Authorization Services in app contexts.
func Collect() (Info, error) {
	raw := C.cmdry_wifi_info()
	if raw == nil {
		return Info{}, nil
	}
	defer C.free(unsafe.Pointer(raw))
	fields := strings.Split(C.GoString(raw), "\t")
	if len(fields) != 5 || fields[0] == "" {
		return Info{}, nil
	}
	info := Info{Available: true, Interface: fields[0], SSID: fields[1], BSSID: fields[2], Signal: fields[3]}
	return info, nil
}
