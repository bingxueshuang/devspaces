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

// spaceCreateCmd represents the spaceCreate command
var spaceCreateCmd = &cobra.Command{
	Use:       "create",
	Short:     "Create a new DevSpace",
	Long:      `Create a new DevSpace.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		tokenFlag, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		pubkey, err := cmd.Flags().GetString("pubkey")
		if err != nil {
			return err
		}
		server := args[0]

		// input
		if server == "" {
			return errors.New("no server url provided")
		}
		if name == "" {
			return errors.New("no devspace name provided")
		}
		pk := new(core.PKey)
		err = keyio.ReadKey(pk, pubkey, "", true)
		if err != nil {
			return err
		}
		token, err := keyio.ReadFile(tokenFlag, false)
		if err != nil {
			return err
		}

		// core
		client := new(http.Client)
		serverURL, err := url.JoinPath(server, "/space/")
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(map[string]any{
			"name":   name,
			"pubkey": hex.EncodeToString(pk.Bytes()),
		})
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", serverURL, buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+string(token))
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
	spaceCmd.AddCommand(spaceCreateCmd)

	spaceCreateCmd.Flags().StringP("name", "n", "", "name of the devspace")
	spaceCreateCmd.Flags().StringP("pubkey", "p", "", "public key of the devspace")
	_ = spaceCreateCmd.MarkFlagRequired("name")
}
