package cmds

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/evertras/cynomys/pkg/cyn"
	"github.com/evertras/cynomys/pkg/httpserver"
	"github.com/evertras/cynomys/pkg/listener"
	"github.com/evertras/cynomys/pkg/sender"
)

var config struct {
	ListenUDP    []string      `mapstructure:"listen-udp"`
	ListenTCP    []string      `mapstructure:"listen-tcp"`
	SendUDP      []string      `mapstructure:"send-udp"`
	SendTCP      []string      `mapstructure:"send-tcp"`
	SendInterval time.Duration `mapstructure:"send-interval"`
	HTTPServer   struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"http"`
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
	flags.String("http.address", "", "An address:port to host an HTTP server on for realtime data, such as '127.0.0.1:8080'")

	err := viper.BindPFlags(flags)

	if err != nil {
		panic(err)
	}
}

func initConfig() {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetEnvPrefix("CYNOMYS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.ReadInConfig()

	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		instance := cyn.New()

		if config.HTTPServer.Address != "" {
			log.Printf("Hosting on http://%s", config.HTTPServer.Address)

			server := httpserver.NewServer(httpserver.Config{
				Addr: config.HTTPServer.Address,
			})

			instance.AddHTTPServer(server)
		}

		for _, listenOnUDP := range config.ListenUDP {
			addr, err := net.ResolveUDPAddr("udp", listenOnUDP)

			if err != nil {
				return fmt.Errorf("net.ResolveUDPAddr for %q: %w", listenOnUDP, err)
			}

			instance.AddUDPListener(listener.NewUDP(*addr))
		}

		for _, listenOnTCP := range config.ListenTCP {
			addr, err := net.ResolveTCPAddr("tcp", listenOnTCP)

			if err != nil {
				return fmt.Errorf("net.ResolveTCPAddr for %q: %w", listenOnTCP, err)
			}

			instance.AddTCPListener(listener.NewTCP(*addr))
		}

		for _, sendUDPTo := range config.SendUDP {
			addr, err := net.ResolveUDPAddr("udp", sendUDPTo)

			if err != nil {
				return fmt.Errorf("net.ResolveUDPAddr for %q: %w", sendUDPTo, err)
			}

			instance.AddUDPSender(sender.NewUDPSender(*addr, config.SendInterval))
		}

		// We could probably generalize this a bit better, but it's short enough
		// not to care for now.
		for _, sendTCPTo := range config.SendTCP {
			addr, err := net.ResolveTCPAddr("tcp", sendTCPTo)

			if err != nil {
				return fmt.Errorf("net.ResolveTCPAddr for %q: %w", sendTCPTo, err)
			}

			instance.AddTCPSender(sender.NewTCPSender(*addr, config.SendInterval))
		}

		return instance.Run()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
