package cmd

import (
	"argocd-tag-plugin/pkg/generate"
	"argocd-tag-plugin/pkg/manifests"
	"fmt"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate <path>",
	Short: "Generate manifests from templates with tags values",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("<path> argument required to generate manifests")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd, args)
	},
}

func run(cmd *cobra.Command, args []string) error {
	m, err := manifests.ReadManifests(args[0], cmd.InOrStdin())
	if err != nil {
		return err
	}

	if err = generate.ProcessManifests(m); err != nil {
		return err
	}

	return manifests.WriteManifests(cmd.OutOrStdout(), m)
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
