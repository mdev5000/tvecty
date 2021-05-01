package main

import (
	"github.com/mdev5000/tvecty"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var noHtml bool
	compileFile := &cobra.Command{
		Use:     "file [*file-in] [file-out]",
		Aliases: []string{"f"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			in, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			var out io.Writer
			if len(args) > 1 {
				out, err := os.OpenFile(args[1], os.O_WRONLY|os.O_CREATE, 0664)
				if err != nil {
					return err
				}
				defer out.Close()
			} else {
				out = os.Stdout
			}
			if noHtml {
				if _, err := tvecty.ExtractHtml(args[0], out, in); err != nil {
					return err
				}
			} else {
				if err := tvecty.ConvertToVecty(args[0], out, in); err != nil {
					return err
				}
			}
			return nil
		},
	}
	compileFile.Flags().BoolVar(&noHtml, "no-html", false, "Do not convert html (useful for debugging).")

	compile := &cobra.Command{
		Use:     "compile",
		Short:   "",
		Aliases: []string{"c"},
	}
	compile.AddCommand(compileFile)

	rootCmd := &cobra.Command{
		Use:   "tvecty [subcommand]",
		Short: "Generate vecty code from templates",
		Args:  cobra.MinimumNArgs(1),
	}
	rootCmd.AddCommand(compile)

	return rootCmd.Execute()
}
