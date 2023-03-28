/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/bingxueshuang/devspaces/cli/keyio"
	"github.com/bingxueshuang/devspaces/core"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

// spaceRequestCmd represents the spaceRequest command
var spaceRequestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request for collaboration",
	Long: `Request for collaboration.

Invite a user for collaboration on a devspace.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}
		devspace, err := cmd.Flags().GetString("devspace")
		if err != nil {
			return err
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		skFlag, err := cmd.Flags().GetString("secret")
		if err != nil {
			return err
		}
		server := args[0]

		// input
		if server == "" {
			return errors.New("server url not provided")
		}
		sk := new(core.SKey)
		err = keyio.ReadKey(sk, skFlag, "", true)
		if err != nil {
			return err
		}

		// core
		client := new(http.Client)
		serverURL, err := url.JoinPath(server, "/space/", devspace, "request")
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(map[string]any{
			"to":     username,
			"secret": hex.EncodeToString(sk.Bytes()),
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
	spaceCmd.AddCommand(spaceRequestCmd)

	spaceRequestCmd.Flags().StringP("devspace", "d", "", "the devspace on which invite is requested")
	spaceRequestCmd.Flags().StringP("username", "u", "", "username of user to be invited")
	spaceRequestCmd.Flags().StringP("secret", "s", "", "secret key of the devspace")
	_ = spaceRequestCmd.MarkFlagRequired("devspace")
	_ = spaceRequestCmd.MarkFlagRequired("username")
}
