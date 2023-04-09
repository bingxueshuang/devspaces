/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bingxueshuang/devspaces/cli/keyio"
	"github.com/bingxueshuang/devspaces/core"
	"github.com/spf13/cobra"
)

func pkFromSkey(cmd *cobra.Command) (core.EllipticKey, error) {
	// flags
	skFile, err := cmd.Flags().GetString("skey")
	if err != nil {
		return nil, err
	}
	skHex, err := cmd.Flags().GetString("skey-hex")
	if err != nil {
		return nil, err
	}

	// input
	sk := new(core.SKey)
	err = keyio.ReadKey(sk, skFile, skHex, true)
	if err != nil {
		return nil, err
	}

	// core logic
	pk := new(core.PKey)
	err = pk.FromSKey(sk)
	return pk, err
}

func pkFromServer(srv string) (core.EllipticKey, error) {
	// core
	client := new(http.Client)
	serverURL, err := url.JoinPath(srv, "/pubkey")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", serverURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	data := new(Response)
	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// output
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	pk, ok := data.Data.(map[string]any)
	if !ok {
		fmt.Println(data)
		return nil, errors.New("invalid response")
	}
	pkHex, ok := pk["pubkey"].(string)
	if !ok {
		fmt.Println(pk)
		return nil, errors.New("invalid response")
	}
	pkey := new(core.PKeyServer)
	pkBytes, err := hex.DecodeString(pkHex)
	if err != nil {
		return nil, err
	}
	err = pkey.FromBytes(pkBytes)
	return pkey, err
}

// pubkeyCmd represents the pubkey command
var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Generate Public Key from a given Private key",
	Long: `Generate Public Key from a given Private key

Take user private key and output the corresponding public key`,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"http://localhost:5005", "http://localhost:8080", "https://api.devspace.com"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var pk core.EllipticKey
		pkFile, err := cmd.Flags().GetString("pkey")
		if err != nil {
			return err
		}
		switch len(args) {
		case 0:
			pk, err = pkFromSkey(cmd)
		case 1:
			pk, err = pkFromServer(args[0])
		default:
			panic("accepts only one argument")
		}
		if err != nil {
			return err
		}
		// output
		return keyio.WriteString(hex.EncodeToString(pk.Bytes()), pkFile, true)
	},
}

func init() {
	rootCmd.AddCommand(pubkeyCmd)

	pubkeyCmd.Flags().StringP("pkey", "p", "", "file to output public key")
	pubkeyCmd.Flags().StringP("skey", "s", "", "private key file")
	_ = pubkeyCmd.MarkFlagFilename("skey") // error happens only when flag does not exist
	pubkeyCmd.Flags().String("skey-hex", "", "hexadecimal user private key")
	pubkeyCmd.MarkFlagsMutuallyExclusive("skey", "skey-hex")
}
