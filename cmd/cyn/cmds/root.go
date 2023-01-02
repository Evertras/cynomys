package cmds

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
	listenOnTCPList []string
	sendUDPToList   []string
	sendTCPToList   []string
	configFilePath  string
	sendInterval    time.Duration
)

func init() {
	rootCmd.Flags().StringSliceVarP(&listenOnUDPList, "listen-udp", "u", nil, "An IP:port address to listen on for UDP.  Can be specified multiple times.")
	rootCmd.Flags().StringSliceVarP(&listenOnTCPList, "listen-tcp", "t", nil, "An IP:port address to listen on for TCP.  Can be specified multiple times.")
	rootCmd.Flags().StringSliceVarP(&sendUDPToList, "send-udp", "U", nil, "An IP:port address to send to (UDP).  Can be specified multiple times.")
	rootCmd.Flags().StringSliceVarP(&sendTCPToList, "send-tcp", "T", nil, "An IP:port address to send to (TCP).  Can be specified multiple times.")
	rootCmd.Flags().StringVarP(&configFilePath, "config-file", "c", "", "A file path to load as additional configuration.")
	rootCmd.Flags().DurationVarP(&sendInterval, "send-interval", "i", time.Second, "How long to wait between attempting to send data")
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

			extraUDPListeners := viper.GetStringSlice("listen-udp")
			listenOnUDPList = append(listenOnUDPList, extraUDPListeners...)

			extraTCPListeners := viper.GetStringSlice("listen-tcp")
			listenOnTCPList = append(listenOnTCPList, extraTCPListeners...)

			extraUDPSenders := viper.GetStringSlice("send-udp")
			sendUDPToList = append(sendUDPToList, extraUDPSenders...)

			extraTCPSenders := viper.GetStringSlice("send-tcp")
			sendTCPToList = append(sendTCPToList, extraTCPSenders...)
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

		for _, listenOnTCP := range listenOnTCPList {
			count++
			addr, err := net.ResolveTCPAddr("tcp", listenOnTCP)

			if err != nil {
				return fmt.Errorf("net.ResolveTCPAddr for %q: %w", listenOnTCP, err)
			}

			eg.Go(func() error {
				l := listener.NewTCP(*addr)

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
					time.Sleep(sendInterval)
				}
			})
		}

		// We could probably generalize this a bit better, but it's short enough
		// not to care for now.
		for _, sendTCPTo := range sendTCPToList {
			count++
			addr, err := net.ResolveTCPAddr("tcp", sendTCPTo)

			if err != nil {
				return fmt.Errorf("net.ResolveTCPAddr for %q: %w", sendTCPTo, err)
			}

			// Shadow capture for use within func below
			sendTCPTo := sendTCPTo

			eg.Go(func() error {
				c := sender.NewTCPSender(*addr)

				for {
					err := c.Send([]byte("hi"))
					if err != nil {
						log.Printf("Failed to send to %q: %v", sendTCPTo, err)
					}
					time.Sleep(sendInterval)
				}
			})
		}

		if count == 0 {
			return fmt.Errorf("no listeners or senders specified")
		}

		return eg.Wait()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
