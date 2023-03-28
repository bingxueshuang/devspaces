/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bingxueshuang/devspaces/cli/keyio"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login a user",
	Long: `Login a user.

Take the username and password of the user and login to
the devspace api server`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		oFlag, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		server := args[0]

		// input
		if server == "" {
			return errors.New("no server url provided")
		}
		if username == "" {
			return errors.New("no username provided")
		}
		if password == "" {
			password, err = keyio.ReadPassword()
			if err != nil {
				return err
			}
		}

		// core logic
		client := http.Client{}
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(map[string]any{
			"username": username,
			"password": password,
		})
		if err != nil {
			return err
		}
		serverURL, err := url.JoinPath(server, "/auth/login")
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", serverURL, buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		data := new(Response)
		err = json.NewDecoder(res.Body).Decode(data)
		if err != nil {
			return err
		}

		// output
		code := res.StatusCode
		if code != http.StatusOK {
			return errors.New(http.StatusText(code))
		}
		datamap, ok := data.Data.(map[string]any)
		if !ok {
			return errors.New("invalid json response")
		}
		token, ok := datamap["token"].(string)
		if !ok {
			return errors.New("invalid json response")
		}
		err = keyio.WriteString(token, oFlag, true)
		return err
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username for Signup")
	loginCmd.Flags().StringP("password", "p", "", "Password for Signup")
	loginCmd.Flags().StringP("output", "o", "", "file to output login token")
	_ = loginCmd.MarkFlagRequired("username")
}
