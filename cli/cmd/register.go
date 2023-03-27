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
	"github.com/bingxueshuang/devspaces/core"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Sign up a user",
	Long: `Sign up a user.

Register a new user to the devspace server using
username, password and public key`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"localhost:5005", "localhost:8080", "api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		pkFlag, err := cmd.Flags().GetString("pkey")
		if err != nil {
			return err
		}
		pkHex, err := cmd.Flags().GetString("pkey-hex")
		if err != nil {
			return err
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		password, err := cmd.Flags().GetString("password")
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
		pk := new(core.PKey)
		err = keyio.ReadKey(pk, pkFlag, pkHex, false)
		if err != nil {
			return err
		}

		// core logic
		client := http.Client{}
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(map[string]any{
			"username": username,
			"password": password,
			"pubkey":   pk,
		})
		if err != nil {
			return err
		}
		serverURL, err := url.JoinPath(server, "/auth/register")
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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("username", "u", "", "Username for Signup")
	registerCmd.Flags().StringP("password", "p", "", "Password for Signup")
	registerCmd.Flags().StringP("pkey", "k", "", "public key file")
	_ = cobra.MarkFlagFilename(registerCmd.Flags(), "pkey")
	registerCmd.Flags().StringP("pkey-hex", "x", "", "hexadecimal public key")
	registerCmd.MarkFlagsMutuallyExclusive("pkey", "pkey-hex")
}
