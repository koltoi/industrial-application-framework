// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package common_types

type AppStatus string

const (
	AppStatusNotSet     = "UNSET"
	AppStatusNotRunning = "NOT_RUNNING"
	AppStatusRunning    = "RUNNING"
	AppStatusFrozen     = "FROZEN"
)

type PrivateNetworkAccess struct {
	Networks             []Network    `json:"networks,omitempty"`
	AppNetwork           string       `json:"appNetwork"`
	NetworkInterfaceName string       `json:"networkInterfaceName,omitempty"`
	AppPodFixIp          *AppPodFixIp `json:"appPodFixIp,omitempty"`
}

type AppPodFixIp struct {
	Db string `json:"db"`
}

type Network struct {
	ApnUUID          string   `json:"apnUUID,omitempty"`
	NetworkID        string   `json:"networkId,omitempty"`
	AdditionalRoutes []string `json:"additionalRoutes,omitempty"`
}