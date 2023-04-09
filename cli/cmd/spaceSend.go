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
	"net/http"
	"net/url"

	"github.com/bingxueshuang/devspaces/cli/keyio"

	"github.com/spf13/cobra"
)

// spaceSendCmd represents the spaceSend command
var spaceSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message on devspace",
	Long: `Send a message on devspace.

This message gets automatically sorted by the server
based on the encrypted keyword using PEKS.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		tokenFlag, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}
		devspace, err := cmd.Flags().GetString("devspace")
		if err != nil {
			return err
		}
		msgFlag, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}
		kwFlag, err := cmd.Flags().GetString("keyword")
		if err != nil {
			return err
		}
		kwHex, err := cmd.Flags().GetString("keyword-hex")
		server := args[0]

		// input
		if server == "" {
			return errors.New("server url not provided")
		}
		msg, err := keyio.ReadFile(msgFlag, true)
		if err != nil {
			return err
		}
		keyword := kwHex
		if keyword == "" {
			kw, err := keyio.ReadFile(kwFlag, false)
			if err != nil {
				return err
			}
			keyword = string(kw)
		}
		token, err := keyio.ReadFile(tokenFlag, false)
		if err != nil {
			return err
		}

		// core
		client := new(http.Client)
		serverURL, err := url.JoinPath(server, "/space/", devspace, "send")
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(map[string]any{
			"data":    hex.EncodeToString(msg),
			"keyword": keyword,
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
	spaceCmd.AddCommand(spaceSendCmd)

	spaceSendCmd.Flags().StringP("devspace", "d", "", "devspace on which message is sent")
	spaceSendCmd.Flags().StringP("message", "m", "", "content of the message")
	spaceSendCmd.Flags().StringP("keyword", "w", "", "encrypted keyword on message")
	spaceSendCmd.Flags().StringP("keyword-hex", "x", "", "hexadecimal keyword")
	spaceSendCmd.MarkFlagsMutuallyExclusive("keyword", "keyword-hex")
	_ = spaceSendCmd.MarkFlagRequired("devspace")
}
