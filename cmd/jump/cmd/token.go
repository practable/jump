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
	"os"
	"time"

	"github.com/ory/viper"
	"github.com/practable/jump/internal/shellaccess"
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "jump token generates a new token for authenticating to jump relay",
	Long: `Set the operating paramters with environment variables, for example

export JUMPTOKEN_LIFETIME=3600
export JUMPTOKEN_ROLE=client
export JUMPTOKEN_SECRET=somesecret
export JUMPTOKEN_TOPIC=123
export JUMPTOKEN_CONNECTIONTYPE=jump
export JUMPTOKEN_AUDIENCE=https://jump-access.example.io
bearer=$(jump token)
`,

	Run: func(cmd *cobra.Command, args []string) {

		viper.SetEnvPrefix("JUMPTOKEN")
		viper.AutomaticEnv()

		lifetime := viper.GetInt64("lifetime")
		role := viper.GetString("role")
		audience := viper.GetString("audience")
		secret := viper.GetString("secret")
		topic := viper.GetString("topic")
		connectionType := viper.GetString("connectionType")

		// check inputs

		if lifetime == 0 {
			fmt.Println("JUMPTOKEN_LIFETIME not set")
			os.Exit(1)
		}
		if role == "" {
			fmt.Println("JUMPTOKEN_ROLE not set")
			os.Exit(1)
		}
		if secret == "" {
			fmt.Println("JUMPTOKEN_SECRET not set")
			os.Exit(1)
		}
		if topic == "" {
			fmt.Println("JUMPTOKEN_TOPIC not set")
			os.Exit(1)
		}
		if connectionType == "" {
			fmt.Println("JUMPTOKEN_CONNECTIONTYPE not set")
			os.Exit(1)
		}
		if audience == "" {
			fmt.Println("JUMPTOKEN_AUDIENCE not set")
			os.Exit(1)
		}

		var scopes []string

		switch role {

		case "host":
			scopes = []string{"host"}
		case "client":
			scopes = []string{"client"}
		case "stats":
			scopes = []string{"stats"}
		case "read":
			scopes = []string{"read"}
		case "write":
			scopes = []string{"write"}
		case "readwrite":
			scopes = []string{"read", "write"}
		default:
			fmt.Println("Unknown role; please choose from host, client, stats, read, write, readwrite")
		}

		iat := time.Now().Unix() - 1 //ensure immediately usable
		nbf := iat
		exp := iat + lifetime

		bearer, err := shellaccess.Token(audience,
			connectionType,
			topic,
			secret,
			scopes,
			iat,
			nbf,
			exp)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(bearer)
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

}
