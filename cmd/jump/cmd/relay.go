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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ory/viper"
	"github.com/practable/jump/internal/shellrelay"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// relayCmd represents the relay command
var relayCmd = &cobra.Command{
	Use:   "jump",
	Short: "jump relay connects jump clients to jump hosts",
	Long: `Set the operating paramters with environment variables, for example


export JUMP_AUDIENCE=https://example.org
export JUMP_BUFFER_SIZE=128
export JUMP_LOG_LEVEL=warn
export JUMP_LOG_FORMAT=json
export JUMP_LOG_FILE=/var/log/relay/relay.log
export JUMP_PORT_ACCESS=3000
export JUMP_PORT_PROFILE=6061
export JUMP_PORT_RELAY=3001
export JUMP_PROFILE=true
export JUMP_SECRET=somesecret
export JUMP_STATS_EVERY=5s
export JUMP_TIDY_EVERY=5m 
export JUMP_URL=wss://example.io/relay 
jump relay

It is expected that you will reverse proxy incoming connections (e.g. with nginx or apache). 
No provision is made for handling TLS in jump relay because this is more convenient
than separately managing certificates, especially when load balancing as may be required.
Note that load balancing takes place at the access phase, with the subsequent connection
being made by the associated relay. The FQDN of your relay access points must be distinct,
so that this affinity is maintained. The FQDN of the access points on the other hand, must
be the same, so that load balancing can be applied in your reverse proxy. All connections 
are individual so connections to a particular host can be simultaneously made in different 
jump relay instances, so long as the relay connection is made in the same instance which
handled the access (see comment above on setting target FQDN to be distinct, so that 
websocket connections are reverse proxied to the correct instance).
`,
	Run: func(cmd *cobra.Command, args []string) {

		viper.SetEnvPrefix("JUMP")
		viper.AutomaticEnv()

		viper.SetDefault("audience", "") //"" so we can check it's been provided
		viper.SetDefault("buffer_size", 128)
		viper.SetDefault("log_file", "/var/log/relay/relay.log")
		viper.SetDefault("log_format", "json")
		viper.SetDefault("log_level", "warn")
		viper.SetDefault("port_access", 3002)
		viper.SetDefault("port_relay", 3003)
		viper.SetDefault("profile", "true")
		viper.SetDefault("profile_port", 6062)
		viper.SetDefault("secret", "") //so we can check it's been provided
		viper.SetDefault("stats_every", "5s")
		viper.SetDefault("tidy_every", "5m")
		viper.SetDefault("url", "") //so we can check it's been provided

		audience := viper.GetString("audience")
		bufferSize := viper.GetInt64("buffer_size")
		logFile := viper.GetString("log_file")
		logFormat := viper.GetString("log_format")
		logLevel := viper.GetString("log_level")
		portAccess := viper.GetInt("port_access")
		portProfile := viper.GetInt("port_profile")
		portRelay := viper.GetInt("port_relay")
		profile := viper.GetBool("profile")
		secret := viper.GetString("secret")
		statsEveryStr := viper.GetString("stats_every")
		tidyEveryStr := viper.GetString("tidy_every")
		URL := viper.GetString("url")

		// Sanity checks
		ok := true

		if audience == "" {
			fmt.Println("You must set JUMP_AUDIENCE")
			ok = false
		}

		if secret == "" {
			fmt.Println("You must set JUMP_SECRET")
			ok = false
		}

		if URL == "" {
			fmt.Println("You must set JUMP_URL")
			ok = false
		}

		if !ok {
			os.Exit(1)
		}

		// parse durations
		statsEvery, err := time.ParseDuration(statsEveryStr)

		if err != nil {
			fmt.Print("cannot parse duration in JUMP_STATS_EVERY=" + statsEveryStr)
			os.Exit(1)
		}

		tidyEvery, err := time.ParseDuration(tidyEveryStr)

		if err != nil {
			fmt.Print("cannot parse duration in JUMP_TIDY_EVERY=" + tidyEveryStr)
			os.Exit(1)
		}

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

		// Report useful info
		log.Infof("jump version: %s", versionString())
		log.Infof("Audience: [%s]", audience)
		log.Infof("Buffer Size: [%d]", bufferSize)
		log.Infof("Log file: [%s]", logFile)
		log.Infof("Log format: [%s]", logFormat)
		log.Infof("Log level: [%s]", logLevel)
		log.Infof("Port for access: [%d]", portAccess)
		log.Infof("Port for profile: [%d]", portProfile)
		log.Infof("Port for relay: [%d]", portRelay)
		log.Infof("Profiling is on: [%t]", profile)
		log.Debugf("Secret: [%s...%s]", secret[:4], secret[len(secret)-4:])
		log.Infof("Stats every: [%s]", statsEvery)
		log.Infof("Tidy every: [%s]", tidyEvery)
		log.Infof("URL: [%s]", URL)

		// Optionally start the profiling server
		if profile {
			go func() {
				url := "localhost:" + strconv.Itoa(portProfile)
				err := http.ListenAndServe(url, nil)
				if err != nil {
					log.Errorf(err.Error())
				}
			}()
		}

		var wg sync.WaitGroup

		closed := make(chan struct{})

		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt)

		go func() {
			for range c {
				close(closed)
				wg.Wait()
				os.Exit(0)
			}
		}()

		wg.Add(1)

		config := shellrelay.Config{
			AccessPort: portAccess,
			Audience:   audience,
			BufferSize: bufferSize,
			RelayPort:  portRelay,
			Secret:     secret,
			StatsEvery: statsEvery,
			Target:     URL,
		}

		go shellrelay.Relay(closed, &wg, config)

		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(relayCmd)
}
