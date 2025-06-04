// Package 帮助函数包
package helpers

import (
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// MicrosecondsStr 将 time.Duration 类型（nano seconds 为单位）
// 输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func MicrosecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func GetDefaultGatewayMac() (mac string, err error) {
	cmd := exec.Command("ip", "route", "show", "default")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	ifaceName := strings.Fields(string(output))[4]
	// 获取当前Mac地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range interfaces {
		if iface.Name == ifaceName {
			mac = iface.HardwareAddr
			break
		}
	}

	return mac, nil
}
