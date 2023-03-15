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
	"strconv"

	"github.com/ory/viper"
	"github.com/practable/jump/internal/shellhost"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "jump host connects a local jump login service to jump relay",
	Long: `Set the operating paramters with environment variables, for example
export JUMP_HOST_LOCAL_PORT=22
export JUMP_HOST_ACCESS=https://access.example.io/jump/shell/abc123
export JUMP_HOST_TOKEN=ey...<snip>
export JUMP_HOST_DEVELOPMENT=true
jump host
`,

	Run: func(cmd *cobra.Command, args []string) {

		viper.SetEnvPrefix("JUMP_HOST")
		viper.AutomaticEnv()

		viper.SetDefault("local_port", 22)
		localPort := viper.GetInt("local_port")
		access := viper.GetString("access")
		token := viper.GetString("token")
		development := viper.GetBool("development")

		if development {
			// development environment
			fmt.Println("Development mode - logging output to stdout")
			fmt.Printf("Local port: %d for %s with %d-byte token\n", localPort, access, len(token))
			log.SetReportCaller(true)
			log.SetFormatter(&log.TextFormatter{})
			log.SetLevel(log.InfoLevel)
			log.SetOutput(os.Stdout)

		} else {

			//production environment
			log.SetFormatter(&log.JSONFormatter{})
			log.SetLevel(log.WarnLevel)
		}

		// check inputs

		if access == "" {
			fmt.Println("JUMP_HOST_ACCESS not set")
			os.Exit(1)
		}
		if token == "" {
			fmt.Println("JUMP_HOST_TOKEN not set")
			os.Exit(1)
		}

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

		local := "localhost:" + strconv.Itoa(localPort)

		go shellhost.Host(ctx, local, access, token)

		<-ctx.Done() //unlikely to exit this way, but block all the same
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}
