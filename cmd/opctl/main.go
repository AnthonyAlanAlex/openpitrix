// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"openpitrix.io/openpitrix/pkg/util/stringutil"
	"openpitrix.io/openpitrix/test"
	"openpitrix.io/openpitrix/test/config"
)

type Flag struct {
	*flag.FlagSet
}

type Cmd interface {
	GetActionName() string
	ParseFlag(f Flag)
	Run(out Out) error
}

// TODO: refactor http client config
var clientConfig = &test.ClientConfig{
	Debug: false,
}

func newRootCmd(c string, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          c,
		Short:        "OpenPitrix cli tool",
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()

	cmd.AddCommand(getCobraCmds(AllCmd)...)
	cmd.AddCommand(getValidateCmd())
	cmd.AddCommand(getCompletionCmd())
	flags.Parse(args)
	return cmd
}

func getCobraCmds(cmds []Cmd) (cobraCmds []*cobra.Command) {
	for _, cmd := range cmds {
		action := cmd.GetActionName()
		underscoreAction := stringutil.CamelCaseToUnderscore(action)
		run := cmd.Run
		c := &cobra.Command{
			Use: fmt.Sprintf("%s [flags]", underscoreAction),
			RunE: func(c *cobra.Command, args []string) error {
				return run(Out{
					action: action,
					out:    c.OutOrStdout(),
				})
			},
		}
		f := c.Flags()
		config.AddFlag(f, &clientConfig.ConfigPath)
		cmd.ParseFlag(Flag{f})

		cobraCmds = append(cobraCmds, c)
	}
	return
}

func main() {
	cmd := newRootCmd(os.Args[0], os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
