package cid

import (
	"fmt"
	"hash"
	"sort"
	"strconv"
)

// RuntimeParams represents basic parameters for
// execution, networking, and volumes for a container.
type RuntimeParams struct {
	User             string            `json:"user"`
	Group            string            `json:"group"`
	CPUShares        uint64            `json:"cpuShares"`
	Memory           uint64            `json:"memory"`
	MemorySwap       uint64            `json:"memorySwap"`
	WorkingDirectory string            `json:"workingDirectory"`
	Ports            []PortSpec        `json:"ports"`
	Volumes          []string          `json:"volumes"`
	Entrypoint       []string          `json:"entrypoint"`
	Command          []string          `json:"command"`
	Environment      map[string]string `json:"environment"`
}

// PortSpec represents a port which a container
// runtime should expose to a containers network.
type PortSpec struct {
	Protocol string `json:"protocol"`
	Port     string `json:"port"`
}

func (rp *RuntimeParams) hash(hasher hash.Hash) {
	// Compute the number of strings we will be hashing.
	length := 6 + len(rp.Volumes) + len(rp.Ports)
	length += len(rp.Entrypoint) + len(rp.Command) + 2*len(rp.Environment)

	// Make a slice that has sufficient capacity.
	values := make([]string, 0, length)

	// Add the first 6 simple values.
	values = append(values,
		rp.User, rp.Group,
		strconv.FormatUint(rp.CPUShares, 10),
		strconv.FormatUint(rp.Memory, 10),
		strconv.FormatUint(rp.MemorySwap, 10),
		rp.WorkingDirectory,
	)

	// Normalize ports as <portnum>/<protocol> and ensure
	// they are in sorted order without duplicates.
	portMap := make(map[string]struct{}, len(rp.Ports))
	for _, port := range rp.Ports {
		key := fmt.Sprintf("%s/%s", port.Port, port.Protocol)
		portMap[key] = struct{}{}
	}
	ports := make([]string, 0, len(portMap))
	for key := range portMap {
		ports = append(ports, key)
	}

	// Sort ports.
	sort.Strings(ports)

	// Sort volumes.
	sort.Strings(rp.Volumes)

	// Add port, volumes, entrypoint, and command.
	values = append(values, ports...)
	values = append(values, rp.Volumes...)
	values = append(values, rp.Entrypoint...)
	values = append(values, rp.Command...)

	// Sort environment variables.
	envKeys := make([]string, 0, len(rp.Environment))
	for key := range rp.Environment {
		envKeys = append(envKeys, key)
	}
	sort.Strings(envKeys)

	// Add environment vars in sorted order.
	for _, key := range envKeys {
		values = append(values, key, rp.Environment[key])
	}

	// Finally, write all of the values to the hasher.
	for _, value := range values {
		hasher.Write([]byte(value))
	}
}
