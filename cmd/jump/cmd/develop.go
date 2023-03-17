/*
Copyright © 2020 Tim Drysdale <timothy.d.drysdale@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/ory/viper"
	"github.com/practable/jump/internal/pipe"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// packetCmd represents the packet command
var developCmd = &cobra.Command{
	Use:   "develop",
	Short: "develop connects two ports using the latest implementation under test",
	Long: `Set the operating paramters with environment variables, for example
export JUMP_PIPE_DEVELOP_PORT_LISTEN=2222
export JUMP_PIPE_DEVELOP_PORT_TARGET=22
export JUMP_PIPE_DEVELOP_LOG_LEVEL=warn
export JUMP_PIPE_DEVELOP_LOG_FORMAT=json
jump pipe packet
`,

	Run: func(cmd *cobra.Command, args []string) {

		viper.SetEnvPrefix("JUMP_PIPE_DEVELOP")
		viper.AutomaticEnv()

		viper.SetDefault("log_format", "json")
		viper.SetDefault("log_level", "warn")
		viper.SetDefault("port_listen", 2222)
		viper.SetDefault("port_target", 22)

		logFormat := viper.GetString("log_format")
		logLevel := viper.GetString("log_level")
		portListen := viper.GetInt("port_listen")
		portTarget := viper.GetInt("port_target")

		// set up logging
		switch strings.ToLower(logLevel) {
		case "trace":
			log.SetLevel(log.TraceLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		case "panic":
			log.SetLevel(log.PanicLevel)
		default:
			fmt.Println("log level can be trace, debug, info, warn, error, fatal or panic but not " + logLevel)
			os.Exit(1)
		}

		switch strings.ToLower(logFormat) {
		case "json":
			log.SetFormatter(&log.JSONFormatter{})
		case "text":
			log.SetFormatter(&log.TextFormatter{})
		default:
			fmt.Println("log format can be json or text but not " + logLevel)
			os.Exit(1)
		}

		// short lived diagnostic utility so log to stdout only
		log.SetOutput(os.Stdout)

		// Report useful info
		log.Infof("jump version: %s", versionString())
		log.Infof("Log format: [%s]", logFormat)
		log.Infof("Log level: [%s]", logLevel)
		log.Infof("Port for listening: [%d]", portListen)
		log.Infof("Port for target: [%d]", portTarget)

		ctx, cancel := context.WithCancel(context.Background())

		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt)

		go func() {
			for range c {
				cancel()
				return
			}
		}()

		config := pipe.Config{
			Listen: portListen,
			Target: portTarget,
		}

		p := pipe.New(config)

		p.RunDevelop(ctx)

	},
}

func init() {
	pipeCmd.AddCommand(developCmd)
}
