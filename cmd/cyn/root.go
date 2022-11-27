package main

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/evertras/cynomys/pkg/listener"
)

var (
	listenOnUDPList []string
)

func init() {
	rootCmd.Flags().StringSliceVarP(&listenOnUDPList, "listen-udp", "l", nil, "An IP:port address to listen on for UDP.  Can be specified multiple times.")
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		eg := errgroup.Group{}

		for _, listenOnUDP := range listenOnUDPList {
			addr, err := net.ResolveUDPAddr("udp", listenOnUDP)

			if err != nil {
				return fmt.Errorf("net.ResolveUDPAddr for %q: %w", listenOnUDP, err)
			}

			eg.Go(func() error {
				l := listener.NewUDP(*addr)

				return l.Listen()
			})
		}

		return eg.Wait()
	},
}
