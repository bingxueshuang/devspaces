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

// tagsListCmd represents the tagsList command
var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags under devspace",
	Long: `List tags under devspace.

Given a devspace, fetch and list the tags under it.`,
	Args:             cobra.ExactArgs(1),
	ValidArgs:        []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		devspace, err := tagsCmd.PersistentFlags().GetString("devspace")
		if err != nil {
			return err
		}
		token, err := tagsCmd.PersistentFlags().GetString("token")
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
		serverURL, err := url.JoinPath(server, "/space/", devspace)
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
		tags, ok := data.Data.([]any)
		if !ok {
			return errors.New("invalid json response")
		}
		err = json.NewEncoder(cmd.OutOrStdout()).Encode(tags)
		return err
	},
}

func init() {
	tagsCmd.AddCommand(tagsListCmd)
}
