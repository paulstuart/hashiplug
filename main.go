// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/paulstuart/hashiplug/shared"
)

var (
	PluginFile = "./kv-plugin"
	verbose    bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "show logs")
	flag.StringVar(&PluginFile, "plugin", PluginFile, "plugin to execute")
}

func main() {
	flag.Parse()

	if !verbose {
		// We don't want to see the plugin logs.
		log.SetOutput(io.Discard)
	}

	plugins := map[int]plugin.PluginSet{}

	// Both version can be supported, but switch the implementation to
	// demonstrate version negoation.
	switch os.Getenv("KV_PROTO") {
	case "netrpc":
		plugins[2] = plugin.PluginSet{
			shared.PluginName: &shared.KVPlugin{},
		}
	case "grpc":
		plugins[3] = plugin.PluginSet{
			shared.PluginName: &shared.KVGRPCPlugin{},
		}
	default:
		fmt.Println("must set KV_PROTO to netrpc or grpc")
		os.Exit(1)
	}

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		VersionedPlugins: plugins,
		Cmd:              exec.Command(PluginFile),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(shared.PluginName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	kv := raw.(shared.KV)
	args := flag.Args()
	switch args[0] {
	case "get":
		result, err := kv.Get(args[1])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		fmt.Println(string(result))

	case "put":
		err := kv.Put(args[1], []byte(args[2]))
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

	case "keys":
		keys, err := kv.Keys()
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		fmt.Printf("keys: %v\n", keys)

	default:
		fmt.Println("Please only use 'get' or 'put'")
		os.Exit(1)
	}
}
