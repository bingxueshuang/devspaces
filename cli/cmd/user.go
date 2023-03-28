/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"encoding/json"
	"errors"
	"github.com/bingxueshuang/devspaces/cli/keyio"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Fetch user public key",
	Long: `Fetch user public key.

Request the devspace api server for public key of
the user`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		oFlag, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		server := args[0]

		// input
		if username == "" {
			return errors.New("no username provided")
		}
		if server == "" {
			return errors.New("no server provided")
		}

		// core
		client := new(http.Client)
		serverURL, err := url.JoinPath(server, "/user/", username)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("GET", serverURL, nil)
		if err != nil {
			return err
		}
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
		if res.StatusCode != http.StatusOK {
			return errors.New(res.Status)
		}
		datamap, ok := data.Data.(map[string]any)
		if !ok {
			return errors.New("invalid json response")
		}
		pubkey, ok := datamap["pubkey"].(string)
		if !ok {
			return errors.New("invalid json response")
		}
		err = keyio.WriteFile([]byte(pubkey), oFlag, true)
		return err
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.Flags().StringP("username", "u", "", "Requested username")
	userCmd.Flags().StringP("output", "o", "", "file to output pubkey")
	_ = userCmd.MarkFlagRequired("username")
}
