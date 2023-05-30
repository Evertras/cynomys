package cmds

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/evertras/cynomys/pkg/listener"
	"github.com/evertras/cynomys/pkg/sender"
)

var config struct {
	ListenUDP    []string      `mapstructure:"listen-udp"`
	ListenTCP    []string      `mapstructure:"listen-tcp"`
	SendUDP      []string      `mapstructure:"send-udp"`
	SendTCP      []string      `mapstructure:"send-tcp"`
	SendInterval time.Duration `mapstructure:"send-interval"`
}

var (
	configFilePath string
)

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.Flags()

	// Special flag for config
	flags.StringVarP(&configFilePath, "config-file", "c", "", "A file path to load as additional configuration.")

	flags.StringSliceP("listen-udp", "u", nil, "An IP:port address to listen on for UDP.  Can be specified multiple times.")
	flags.StringSliceP("listen-tcp", "t", nil, "An IP:port address to listen on for TCP.  Can be specified multiple times.")
	flags.StringSliceP("send-udp", "U", nil, "An IP:port address to send to (UDP).  Can be specified multiple times.")
	flags.StringSliceP("send-tcp", "T", nil, "An IP:port address to send to (TCP).  Can be specified multiple times.")
	flags.DurationP("send-interval", "i", time.Second, "How long to wait between attempting to send data")

	viper.BindPFlags(flags)
}

func initConfig() {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	}

	viper.ReadInConfig()

	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		eg := errgroup.Group{}
		count := 0

		for _, listenOnUDP := range config.ListenUDP {
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

		for _, listenOnTCP := range config.ListenTCP {
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

		for _, sendUDPTo := range config.SendUDP {
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
					time.Sleep(config.SendInterval)
				}
			})
		}

		// We could probably generalize this a bit better, but it's short enough
		// not to care for now.
		for _, sendTCPTo := range config.SendTCP {
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
					time.Sleep(config.SendInterval)
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
