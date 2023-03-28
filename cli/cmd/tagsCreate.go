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

// tagsCreateCmd represents the tagsCreate command
var tagsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create tags under devspace",
	Long: `Create tags under devspace.

Given a particular devspace and if the user has
permission to create tags on it (the owner), then
create a new tag and add it under the devspace.`,
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
		tdFlag, err := cmd.Flags().GetString("trapdoor")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		server := args[0]

		// input
		tdBytes, err := keyio.ReadFile(tdFlag, true)
		if err != nil {
			return err
		}
		trapdoor := string(tdBytes)
		if server == "" {
			return errors.New("no server url supplied")
		}
		if name == "" {
			return errors.New("no tag name supplied")
		}

		// core
		client := new(http.Client)
		serverURL, err := url.JoinPath(server, "/space/", devspace)
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(map[string]any{
			"name":     name,
			"trapdoor": trapdoor,
		})
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", serverURL, buf)
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
		return nil
	},
}

func init() {
	tagsCmd.AddCommand(tagsCreateCmd)

	tagsCreateCmd.Flags().StringP("name", "n", "", "name of the tag")
	tagsCreateCmd.Flags().StringP("trapdoor", "t", "", "trapdoor for the tag")
	_ = tagsCreateCmd.MarkFlagRequired("name")
}
