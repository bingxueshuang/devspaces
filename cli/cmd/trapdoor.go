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

// trapdoorCmd represents the trapdoor command
var trapdoorCmd = &cobra.Command{
	Use:   "trapdoor",
	Short: "Get trapdoor for a given keyword",
	Long: `Get PEKS trapdoor for a given keyword.

Take the keyword, server public key, receiver secret key and
sender public key from the command line and output trapdoor
for the given keyword`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error { // flags
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
		peks, err := core.Trapdoor([]byte(keyword), server, pk, sk)
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
	rootCmd.AddCommand(trapdoorCmd)

	trapdoorCmd.Flags().StringP("output", "o", "", "file to output trapdoor")
	trapdoorCmd.Flags().StringP("skey", "s", "", "private key file")
	_ = trapdoorCmd.MarkFlagFilename("skey")
	trapdoorCmd.Flags().StringP("pkey", "p", "", "public key file")
	_ = trapdoorCmd.MarkFlagFilename("pkey")
	trapdoorCmd.Flags().StringP("server", "r", "", "server public key file")
	_ = trapdoorCmd.MarkFlagFilename("server")
	trapdoorCmd.Flags().String("server-hex", "", "hexadecimal server public key")
	trapdoorCmd.Flags().String("skey-hex", "", "hexadecimal private key")
	trapdoorCmd.Flags().String("pkey-hex", "", "hexadecimal public key")
	trapdoorCmd.Flags().StringP("file", "f", "", "keyword file")
	_ = trapdoorCmd.MarkFlagFilename("file")
	trapdoorCmd.Flags().StringP("keyword", "k", "", "keyword text")
	trapdoorCmd.MarkFlagsMutuallyExclusive("skey", "skey-hex")
	trapdoorCmd.MarkFlagsMutuallyExclusive("pkey", "pkey-hex")
	trapdoorCmd.MarkFlagsMutuallyExclusive("file", "keyword")
}
