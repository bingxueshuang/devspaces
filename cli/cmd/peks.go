/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"encoding/hex"
	"github.com/bingxueshuang/devspaces/cli/keyio"
	"github.com/bingxueshuang/devspaces/core"
	"github.com/spf13/cobra"
)

// peksCmd represents the peks command
var peksCmd = &cobra.Command{
	Use:   "peks",
	Short: "Get Public key searchable encryption of provided keyword",
	Long: `Get Public key searchable encryption of provided keyword.

Take keyword, sender private key, receiver public key and server public key
as input and output the public key searchable encryption of the keyword.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		oFlag, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		srvFlag, err := cmd.Flags().GetString("server")
		if err != nil {
			return err
		}
		skFlag, err := cmd.Flags().GetString("skey")
		if err != nil {
			return err
		}
		pkFlag, err := cmd.Flags().GetString("pkey")
		if err != nil {
			return err
		}
		srvHex, err := cmd.Flags().GetString("server-hex")
		if err != nil {
			return err
		}
		skHex, err := cmd.Flags().GetString("skey-hex")
		if err != nil {
			return err
		}
		pkHex, err := cmd.Flags().GetString("pkey-hex")
		if err != nil {
			return err
		}
		kwFlag, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		keyword, err := cmd.Flags().GetString("keyword")
		if err != nil {
			return err
		}

		// input
		sk := new(core.SKey)
		err = keyio.ReadKey(sk, skFlag, skHex, false)
		if err != nil {
			return err
		}
		pk := new(core.PKey)
		err = keyio.ReadKey(pk, pkFlag, pkHex, false)
		if err != nil {
			return err
		}
		server := new(core.PKeyServer)
		err = keyio.ReadKey(server, srvFlag, srvHex, false)
		if err != nil {
			return err
		}
		if keyword == "" {
			kwBytes, err := keyio.ReadFile(kwFlag, true)
			if err != nil {
				return err
			}
			keyword = string(kwBytes)
		}

		// core logic
		peks, err := core.PEKS([]byte(keyword), server, pk, sk)
		if err != nil {
			return err
		}

		// output
		err = keyio.WriteString(hex.EncodeToString(peks), oFlag, true)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(peksCmd)

	peksCmd.Flags().StringP("output", "o", "", "file to output shared key")
	peksCmd.Flags().StringP("skey", "s", "", "private key file")
	_ = peksCmd.MarkFlagFilename("skey")
	peksCmd.Flags().StringP("pkey", "p", "", "public key file")
	_ = peksCmd.MarkFlagFilename("pkey")
	peksCmd.Flags().StringP("server", "r", "", "server public key file")
	_ = peksCmd.MarkFlagFilename("server")
	peksCmd.Flags().String("server-hex", "", "hexadecimal server public key")
	peksCmd.Flags().String("skey-hex", "", "hexadecimal private key")
	peksCmd.Flags().String("pkey-hex", "", "hexadecimal public key")
	peksCmd.Flags().StringP("file", "f", "", "keyword file")
	_ = peksCmd.MarkFlagFilename("file")
	peksCmd.Flags().StringP("keyword", "k", "", "keyword text")
	peksCmd.MarkFlagsMutuallyExclusive("skey", "skey-hex")
	peksCmd.MarkFlagsMutuallyExclusive("pkey", "pkey-hex")
	peksCmd.MarkFlagsMutuallyExclusive("file", "keyword")
}
