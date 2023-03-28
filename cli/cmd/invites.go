/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

// invitesCmd represents the invites command
var invitesCmd = &cobra.Command{
	Use:       "invites",
	Short:     "List all collaboration invites",
	Long:      `List all collaboration invites.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}
		server := args[0]

		// input
		if server == "" {
			return errors.New("server url not supplied")
		}

		// core
		client := new(http.Client)
		serverURL, err := url.JoinPath(server, "/dashboard/")
		if err != nil {
			return err
		}
		req, err := http.NewRequest("GET", serverURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
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
		devs, ok := data.Data.([]any)
		if !ok {
			return errors.New("invalid json response")
		}
		err = json.NewEncoder(cmd.OutOrStdout()).Encode(devs)
		return err
	},
}

func init() {
	rootCmd.AddCommand(invitesCmd)
	invitesCmd.Flags().StringP("token", "k", "", "login token")
}
