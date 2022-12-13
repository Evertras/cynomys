package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/evertras/cynomys/pkg/listener"
	"github.com/evertras/cynomys/pkg/sender"
)

var (
	listenOnUDPList []string
	sendUDPToList   []string
	configFilePath  string
)

func init() {
	rootCmd.Flags().StringSliceVarP(&listenOnUDPList, "listen-udp", "u", nil, "An IP:port address to listen on for UDP.  Can be specified multiple times.")
	rootCmd.Flags().StringSliceVarP(&sendUDPToList, "send-udp", "U", nil, "An IP:port address to send to (UDP).  Can be specified multiple times.")
	rootCmd.Flags().StringVarP(&configFilePath, "config-file", "c", "", "A file path to load as additional configuration.")
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		if configFilePath != "" {
			file, err := os.Open(configFilePath)

			if err != nil {
				return fmt.Errorf("failed to open config file %q: %w", configFilePath, err)
			}

			viper.SetConfigType("yaml")
			err = viper.ReadConfig(file)

			if err != nil {
				return fmt.Errorf("failed to read config %q: %w", configFilePath, err)
			}

			extraUdpListeners := viper.GetStringSlice("listen-udp")
			listenOnUDPList = append(listenOnUDPList, extraUdpListeners...)

			extraUdpSenders := viper.GetStringSlice("send-udp")
			sendUDPToList = append(sendUDPToList, extraUdpSenders...)
		}

		eg := errgroup.Group{}
		count := 0

		for _, listenOnUDP := range listenOnUDPList {
			count++
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
			count++
			addr, err := net.ResolveUDPAddr("udp", sendUDPTo)

			if err != nil {
				return fmt.Errorf("net.ResolveUDPAddr for %q: %w", sendUDPTo, err)
			}

			// Shadow capture for use within func below
			sendUDPTo := sendUDPTo

			eg.Go(func() error {
				c := sender.NewUDPSender(*addr)

				for {
					err := c.Send([]byte("hi"))
					if err != nil {
						log.Printf("Failed to send to %q: %v", sendUDPTo, err)
					}
					time.Sleep(time.Second)
				}
			})
		}

		if count == 0 {
			return fmt.Errorf("no listeners or senders specified")
		}

		return eg.Wait()
	},
}
