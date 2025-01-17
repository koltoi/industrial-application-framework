// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package consulclient

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/serviceconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/json"
)

func CreateConsulClient() (*consulapi.Client, error) {
	log.Info("CreateConsulClient called")

	conf := consulapi.DefaultConfig()
	conf.Address = fmt.Sprintf(serviceconfig.ConfigData.ConsulAddress)

	consulClient, err := consulapi.NewClient(conf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to creat consul api client")
	}
	return consulClient, nil
}

func ReadConsulContent(consulClient *consulapi.Client) (string, error) {
	log.Info("ReadConsulContent called")

	KVPairs, _, err := consulClient.KV().List("/", nil)
	if err != nil {
		return "", errors.Wrap(err, "Failed to list consul content")
	}
	log.Info("consul content", "KVPairs", KVPairs)

	consulContent, err := json.Marshal(KVPairs)
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal the KVPairs map")
	}

	return string(consulContent), nil
}

func AddEntryToConsul(consulClient *consulapi.Client, key, value string) (error) {
	log.Info("AddEntryToConsul called")

	d := &consulapi.KVPair{Key: key, Value: []byte(value)}
	_, err := consulClient.KV().Put(d, nil)

	if err != nil {
		errors.Wrap(err, "Failed to add entry to consul")
	}
	return err
}

