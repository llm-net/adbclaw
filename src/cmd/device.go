package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// DeviceEntry represents a connected device from "adb devices -l".
type DeviceEntry struct {
	Serial      string `json:"serial"`
	State       string `json:"state"`
	Product     string `json:"product,omitempty"`
	Model       string `json:"model,omitempty"`
	Device      string `json:"device,omitempty"`
	TransportID string `json:"transport_id,omitempty"`
}

// DeviceInfo holds device properties from getprop.
type DeviceInfo struct {
	Serial       string `json:"serial,omitempty"`
	Model        string `json:"model"`
	Brand        string `json:"brand"`
	Manufacturer string `json:"manufacturer"`
	AndroidVer   string `json:"android_version"`
	SDKLevel     string `json:"sdk_level"`
	ScreenSize   string `json:"screen_size,omitempty"`
	Density      string `json:"density,omitempty"`
	ABIs         string `json:"abis"`
}

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Device management commands",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		result, err := client.RawCommand("devices", "-l")
		if err != nil {
			writer.Fail("device list", "ADB_ERROR", err.Error(),
				"Check that adb is installed and accessible", start)
			return nil
		}

		devices := parseDeviceList(result.Stdout)
		writer.Success("device list", map[string]interface{}{
			"devices": devices,
			"count":   len(devices),
		}, start)
		return nil
	},
}

var deviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show device details",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		info, err := getDeviceInfo()
		if err != nil {
			writer.Fail("device info", "ADB_ERROR", err.Error(),
				"Ensure a device is connected: adbclaw device list", start)
			return nil
		}
		writer.Success("device info", info, start)
		return nil
	},
}

func init() {
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceInfoCmd)
	rootCmd.AddCommand(deviceCmd)
}

func parseDeviceList(output string) []DeviceEntry {
	var devices []DeviceEntry
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		entry := DeviceEntry{
			Serial: parts[0],
			State:  parts[1],
		}
		// Parse key:value pairs like "product:xxx model:yyy"
		for _, p := range parts[2:] {
			kv := strings.SplitN(p, ":", 2)
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "product":
				entry.Product = kv[1]
			case "model":
				entry.Model = kv[1]
			case "device":
				entry.Device = kv[1]
			case "transport_id":
				entry.TransportID = kv[1]
			}
		}
		devices = append(devices, entry)
	}
	return devices
}

func getDeviceInfo() (*DeviceInfo, error) {
	// Single getprop call to get all properties at once (instead of 6 separate calls).
	// Output format: [key]: [value]
	result, err := client.Shell("getprop")
	if err != nil {
		return nil, err
	}
	props := parseGetprop(result.Stdout)

	// Get screen size and density in parallel with a combined wm call
	screenSize := ""
	density := ""

	type wmResult struct {
		size    string
		density string
	}
	wmCh := make(chan wmResult, 1)
	go func() {
		var r wmResult
		if res, err := client.Shell("wm", "size"); err == nil {
			for _, line := range strings.Split(res.Stdout, "\n") {
				if strings.Contains(line, "Physical size") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						r.size = strings.TrimSpace(parts[1])
					}
				}
			}
		}
		if res, err := client.Shell("wm", "density"); err == nil {
			for _, line := range strings.Split(res.Stdout, "\n") {
				if strings.Contains(line, "Physical density") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						r.density = strings.TrimSpace(parts[1])
					}
				}
			}
		}
		wmCh <- r
	}()
	wm := <-wmCh
	screenSize = wm.size
	density = wm.density

	info := &DeviceInfo{
		Serial:       flagSerial,
		Model:        props["ro.product.model"],
		Brand:        props["ro.product.brand"],
		Manufacturer: props["ro.product.manufacturer"],
		AndroidVer:   props["ro.build.version.release"],
		SDKLevel:     props["ro.build.version.sdk"],
		ABIs:         props["ro.product.cpu.abilist"],
		ScreenSize:   screenSize,
		Density:      density,
	}
	return info, nil
}

// parseGetprop parses the output of "adb shell getprop" into a map.
// Each line has the format: [key]: [value]
func parseGetprop(output string) map[string]string {
	props := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "[") {
			continue
		}
		// Format: [key]: [value]
		closeBracket := strings.Index(line, "]")
		if closeBracket < 2 {
			continue
		}
		key := line[1:closeBracket]
		// Skip ": [" separator
		rest := line[closeBracket+1:]
		valStart := strings.Index(rest, "[")
		valEnd := strings.LastIndex(rest, "]")
		if valStart < 0 || valEnd <= valStart {
			continue
		}
		value := rest[valStart+1 : valEnd]
		props[key] = value
	}
	return props
}
