// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/paulstuart/hashiplug/shared"
	"github.com/paulstuart/hashiplug/shared/exec"
)

// Here is a real implementation of KV that uses grpc and  writes to a local
// file with the key name and the contents are the value of the key.
/*
type KVGRPC struct{}

func (KVGRPC) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from plugin version 3\n", string(value)))
	return os.WriteFile(shared.FilePrefix+key, value, 0644)
}

func (KVGRPC) Get(key string) ([]byte, error) {
	d, err := os.ReadFile(shared.FilePrefix + key)
	if err != nil {
		return nil, err
	}
	return append(d, []byte("Read by plugin version 3\n")...), nil
}

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type KV struct{}

func (KV) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from plugin version 2\n", string(value)))
	return os.WriteFile(shared.FilePrefix+key, value, 0644)
}

func (KV) Get(key string) ([]byte, error) {
	d, err := os.ReadFile(shared.FilePrefix + key)
	if err != nil {
		return nil, err
	}
	return append(d, []byte("Read by plugin version 2\n")...), nil
}
*/

type KVGRPC = exec.KVGRPC
type KV = exec.KV

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			// Version 2 only uses NetRPC
			2: {
				shared.PluginName: &shared.KVPlugin{Impl: &KV{}},
			},
			// Version 3 only uses GRPC
			3: {
				shared.PluginName: &shared.KVGRPCPlugin{Impl: &KVGRPC{}},
			},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
