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
	"github.com/evertras/cynomys/pkg/listener"
	"github.com/evertras/cynomys/pkg/metrics"
	"github.com/evertras/cynomys/pkg/sender"
)

var config struct {
	Listen struct {
		Udp []string `mapstructure:"udp"`
		Tcp []string `mapstructure:"tcp"`
	} `mapstructure:"listen"`

	Send struct {
		Udp []string `mapstructure:"udp"`
		Tcp []string `mapstructure:"tcp"`

		Interval time.Duration `mapstructure:"interval"`
		Data     string        `mapstructure:"data"`
	} `mapstructure:"send"`

	HTTPServer struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"http"`

	Sinks struct {
		SinkStdout struct {
			Enabled bool `mapstructure:"enabled"`
		} `mapstructure:"stdout"`
	} `mapstructure:"sinks"`
}

var (
	configFilePath string
)

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.Flags()

	// Special flag for config
	flags.StringVarP(&configFilePath, "config", "c", "", "A file path to load as additional configuration.")

	flags.StringSliceP("listen.udp", "u", nil, "An IP:port address to listen on for UDP.  Can be specified multiple times.")
	flags.StringSliceP("listen.tcp", "t", nil, "An IP:port address to listen on for TCP.  Can be specified multiple times.")
	flags.StringSliceP("send.udp", "U", nil, "An IP:port address to send to (UDP).  Can be specified multiple times.")
	flags.StringSliceP("send.tcp", "T", nil, "An IP:port address to send to (TCP).  Can be specified multiple times.")
	flags.StringP("send.data", "d", "hi", "The string data to send.")
	flags.DurationP("send.interval", "i", time.Second, "How long to wait between attempting to send data")
	flags.String("http.address", "", "An address:port to host an HTTP server on for realtime data, such as '127.0.0.1:8080'")
	flags.Bool("sinks.stdout.enabled", false, "Whether to enable the stdout metrics sink")

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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Ignore errors here because we don't necessarily need a config file
	_ = viper.ReadInConfig()

	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		sinks := []metrics.Sink{}

		if config.Sinks.SinkStdout.Enabled {
			sinks = append(sinks, metrics.NewSinkStdout())
		}

		sink := metrics.NewMultiSink(sinks...)

		instance := cyn.New()

		if config.HTTPServer.Address != "" {
			log.Printf("Hosting on http://%s", config.HTTPServer.Address)

			instance.AddHTTPServer(config.HTTPServer.Address)
		}

		for _, listenOnUDP := range config.Listen.Udp {
			addr, err := net.ResolveUDPAddr("udp", listenOnUDP)

			if err != nil {
				return fmt.Errorf("net.ResolveUDPAddr for %q: %w", listenOnUDP, err)
			}

			instance.AddUDPListener(listener.NewUDP(*addr))
		}

		for _, listenOnTCP := range config.Listen.Tcp {
			addr, err := net.ResolveTCPAddr("tcp", listenOnTCP)

			if err != nil {
				return fmt.Errorf("net.ResolveTCPAddr for %q: %w", listenOnTCP, err)
			}

			instance.AddTCPListener(listener.NewTCP(*addr))
		}

		for _, sendUDPTo := range config.Send.Udp {
			addr, err := net.ResolveUDPAddr("udp", sendUDPTo)

			if err != nil {
				return fmt.Errorf("net.ResolveUDPAddr for %q: %w", sendUDPTo, err)
			}

			instance.AddUDPSender(sender.NewUDPSender(*addr, config.Send.Interval, sink, []byte(config.Send.Data)))
		}

		// We could probably generalize this a bit better, but it's short enough
		// not to care for now.
		for _, sendTCPTo := range config.Send.Tcp {
			addr, err := net.ResolveTCPAddr("tcp", sendTCPTo)

			if err != nil {
				return fmt.Errorf("net.ResolveTCPAddr for %q: %w", sendTCPTo, err)
			}

			instance.AddTCPSender(sender.NewTCPSender(*addr, config.Send.Interval, sink, []byte(config.Send.Data)))
		}

		return instance.Run()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
