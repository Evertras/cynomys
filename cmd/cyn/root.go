package main

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/evertras/cynomys/pkg/listener"
	"github.com/evertras/cynomys/pkg/sender"
)

var (
	listenOnUDPList []string
	sendUDPToList   []string
)

func init() {
	rootCmd.Flags().StringSliceVarP(&listenOnUDPList, "listen-udp", "u", nil, "An IP:port address to listen on for UDP.  Can be specified multiple times.")
	rootCmd.Flags().StringSliceVarP(&sendUDPToList, "send-udp", "U", nil, "An IP:port address to send to (UDP).  Can be specified multiple times.")
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

		for _, sendUDPTo := range sendUDPToList {
			addr, err := net.ResolveUDPAddr("udp", sendUDPTo)

			if err != nil {
				return fmt.Errorf("net.ResolveUDPAddr for %q: %w", sendUDPTo, err)
			}

			eg.Go(func() error {
				c := sender.NewUDPSender(*addr)

				for {
					err := c.Send([]byte("hi"))
					if err != nil {
						return fmt.Errorf("c.Send: %w", err)
					}
					time.Sleep(time.Second)
				}
			})
		}

		return eg.Wait()
	},
}
