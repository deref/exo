package main

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker/components/network"
	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/spf13/cobra"
)

func init() {
	newCmd.AddCommand(newNetworkCmd)
	// TODO: Options from docker to consider adding support for:
	//	--attachable           Enable manual container attachment
	//	--aux-address map      Auxiliary IPv4 or IPv6 addresses used by Network driver (default map[])
	//	--config-from string   The network from which to copy the configuration
	//	--config-only          Create a configuration only network
	//-d, --driver string        Driver to manage the Network (default "bridge")
	//	--gateway strings      IPv4 or IPv6 Gateway for the master subnet
	//	--ingress              Create swarm routing-mesh network
	//	--internal             Restrict external access to the network
	//	--ip-range strings     Allocate container ip from a sub-range
	//	--ipam-driver string   IP Address Management Driver (default "default")
	//	--ipam-opt map         Set IPAM driver specific options (default map[])
	//	--ipv6                 Enable IPv6 networking
	//	--label list           Set metadata on a network
	//-o, --opt map              Set driver specific options (default map[])
	//	--scope string         Control the network's scope
	//	--subnet strings       Subnet in CIDR format that represents a network segment
}

var networkSpec network.Spec

var newNetworkCmd = &cobra.Command{
	Use:   "network <name> [options]",
	Short: "Creates a new network",
	Long: `Creates a new network.

Similar in spirit to:

docker network create <name>
`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)

		name := args[0]

		output, err := workspace.CreateComponent(ctx, &api.CreateComponentInput{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(networkSpec),
		})
		if err != nil {
			return err
		}
		fmt.Println(output.ID)
		return nil
	},
}
