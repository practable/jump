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
	"github.com/practable/jump/internal/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "jump client forwards local jump logins to jump relay",
	Long: `Set the operating paramters with environment variables, for example
export JUMP_CLIENT_LOCAL_PORT=22
export JUMP_CLIENT_RELAY_SESSION=https://access.example.io/jump/abc123
export JUMP_CLIENT_TOKEN=ey...<snip>
export JUMP_CLIENT_LOG_LEVEL=warn
export JUMP_CLIENT_LOG_FORMAT=json
export JUMP_CLIENT_LOG_FILE=/var/log/jump/jump-client.log
jump client
`,
	Run: func(cmd *cobra.Command, args []string) {

		viper.SetEnvPrefix("JUMP_CLIENT")
		viper.AutomaticEnv()

		viper.SetDefault("local_port", 8082)
		viper.SetDefault("log_file", "/var/log/jump/jump-client.log")
		viper.SetDefault("log_format", "json")
		viper.SetDefault("log_level", "warn")

		localPort := viper.GetInt("local_port")
		logFile := viper.GetString("log_file")
		logFormat := viper.GetString("log_format")
		logLevel := viper.GetString("log_level")
		relaySession := viper.GetString("relay_session")
		token := viper.GetString("token")

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
			fmt.Println("BOOK_LOG_LEVEL can be trace, debug, info, warn, error, fatal or panic but not " + logLevel)
			os.Exit(1)
		}
		
		switch strings.ToLower(logFormat) {
		case "json":
			log.SetFormatter(&log.JSONFormatter{})
		case "text":
			log.SetFormatter(&log.TextFormatter{})
		default:
			fmt.Println("BOOK_LOG_FORMAT can be json or text but not " + logLevel)
			os.Exit(1)
		}

		if strings.ToLower(logFile) == "stdout" {

			log.SetOutput(os.Stdout) //

		} else {

			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				log.SetOutput(file)
			} else {
				log.Infof("Failed to log to %s, logging to default stderr", logFile)
			}
		}		



		// check inputs

		if relaySession == "" {
			fmt.Println("JUMP_CLIENT_RELAY_SESSION not set")
			os.Exit(1)
		}
		if token == "" {
			fmt.Println("JUMP_CLIENT_TOKEN not set")
			os.Exit(1)
		}

		log.Infof("jump version: %s", versionString())
		log.Infof("Log file: [%s]", logFile)
		log.Infof("Log format: [%s]", logFormat)
		log.Infof("Log level: [%s]", logLevel)
		log.Infof("Port for local listening: [%d]", localPort)
    	log.Infof("Relay session: [%s]", relaySession)
		log.Infof("Token length: [%d]", len(token))
		
		ctx, cancel := context.WithCancel(context.Background())

		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt)

		go func() {
			for range c {
				cancel()
				<-ctx.Done()
				os.Exit(0)
			}
		}()

		go client.Run(ctx, localPort, relaySession, token)

		<-ctx.Done() //unlikely to exit this way, but block all the same
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
