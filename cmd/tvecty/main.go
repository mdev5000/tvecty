package main

import (
	"fmt"
	"github.com/mdev5000/globerous"
	"github.com/mdev5000/tvecty"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var globCompiler = globerous.NewCompiler(globerous.HybridGlobRegexPartCompiler)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	compile := &cobra.Command{
		Use:     "compile",
		Short:   "",
		Aliases: []string{"c"},
	}
	compile.AddCommand(
		cmdCompileFile(),
		cmdCompileDir(),
	)

	rootCmd := &cobra.Command{
		Use:   "tvecty [subcommand]",
		Short: "Generate vecty code from templates",
		Args:  cobra.MinimumNArgs(1),
	}
	rootCmd.AddCommand(compile)

	return rootCmd.Execute()
}

func cmdCompileDir() *cobra.Command {
	var outSuffix string
	compileFile := cobra.Command{
		Use: "dir [*file-glob]",
		Example: `
  tvecty c d somepackage/*.vtpl
  tvecty c d somepackage/**/*.vtpl`,
		Aliases: []string{"d"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileGlob := args[0]
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			filesToCompile, err := matchFiles(wd, fileGlob, "file-glob")
			if err != nil {

			}
			for _, f := range filesToCompile {
				if err := compileFileToPath(f, f+outSuffix); err != nil {
					return err
				}
			}
			return nil
		},
	}
	compileFile.Flags().StringVarP(&outSuffix, "suffix", "s", ".go", "Suffix to place at the end of the compiled file.")
	return &compileFile
}

func cmdCompileFile() *cobra.Command {
	var noHtml bool
	compileFile := cobra.Command{
		Use:     "file [*file-in] [file-out]",
		Aliases: []string{"f"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var out io.Writer
			if len(args) > 1 {
				out, err := openFile(args[1])
				if err != nil {
					return err
				}
				defer out.Close()
			} else {
				out = os.Stdout
			}
			return compileFile(args[0], out, noHtml)
		},
	}
	compileFile.Flags().BoolVar(&noHtml, "no-html", false, "Do not convert html (useful for debugging).")
	return &compileFile
}

func matchFiles(dir, glob, debugName string) ([]string, error) {
	if strings.HasPrefix(glob, "/") {
		return nil, fmt.Errorf("%s cannot be an absolute path", debugName)
	}
	matchingPathParts := strings.Split(glob, "/")
	m, err := globCompiler.Compile(matchingPathParts...)
	if err != nil {
		return nil, err
	}
	var filesToCompile []string
	err = globerous.WalkSimple(globerous.NewOSGlobFs(), m, dir, func(dir string, info os.FileInfo) error {
		filesToCompile = append(filesToCompile, filepath.Join(dir, info.Name()))
		return nil
	})
	return filesToCompile, err
}

func compileFileToPath(fPathIn, fPathOut string) error {
	out, err := openFile(fPathOut)
	if err != nil {
		return err
	}
	defer out.Close()
	return compileFile(fPathIn, out, false)
}

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664)
}

func compileFile(fPathIn string, out io.Writer, noHtml bool) error {
	in, err := os.ReadFile(fPathIn)
	if err != nil {
		return err
	}
	if noHtml {
		if _, err := tvecty.ExtractHtml(fPathIn, out, in); err != nil {
			return err
		}
	} else {
		if err := tvecty.ConvertToVecty(fPathIn, out, in); err != nil {
			return err
		}
	}
	return nil
}
