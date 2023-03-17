// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"fmt"
	"os"
	"strings"
)

const FilePrefix = "__kv_"

// Here is a real implementation of KV that uses grpc and  writes to a local
// file with the key name and the contents are the value of the key.
type KVGRPC struct{}

func (KVGRPC) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from ENHANCED plugin version 3\n", string(value)))
	return os.WriteFile(FilePrefix+key, value, 0644)
}

func (KVGRPC) Get(key string) ([]byte, error) {
	d, err := os.ReadFile(FilePrefix + key)
	if err != nil {
		return nil, err
	}
	return append(d, []byte("Read by plugin version 3\n")...), nil
}

func (KVGRPC) Keys() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	const pre = len(FilePrefix)
	var keys []string
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, FilePrefix) {
			keys = append(keys, name[pre:])
		}

	}
	return keys, nil
}

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type KV struct{}

func (KV) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from A ROCKIN' plugin version 2\n", string(value)))
	return os.WriteFile(FilePrefix+key, value, 0644)
}

func (KV) Get(key string) ([]byte, error) {
	d, err := os.ReadFile(FilePrefix + key)
	if err != nil {
		return nil, err
	}
	return append(d, []byte("Read by plugin version 2\n")...), nil
}

func (KV) Keys() ([]string, error) {
	return KVGRPC{}.Keys()
}
