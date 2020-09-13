package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/DataReply/korgi/internal/version"

	"github.com/spf13/cobra"
)

const versionDesc = `
Show the version for Korgi.

This will print a representation the version of Korgi.
The output will look something like this:

version.BuildInfo{Version:"v2.0.0", GitCommit:"ff52399e51bb880526e9cd0ed8386f6433b74da1", GitTreeState:"clean"}

- Version is the semantic version of the release.
- GitCommit is the SHA for the commit that this version was built from.
- GitTreeState is "clean" if there are no local code changes when this binary was
  built, and "dirty" if the binary was built from locally modified code.
`

type versionOptions struct {
	short bool
}

func newVersionCmd(out io.Writer) *cobra.Command {
	o := &versionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the client version information",
		Long:  versionDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&o.short, "short", false, "print the version number")
	return cmd
}

func (o *versionOptions) run(out io.Writer) error {
	fmt.Fprintln(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	v := version.Get()
	if short {
		return fmt.Sprintf("%s+%s", v.Version, v.GitCommit[:7])
	}
	return fmt.Sprintf("%#v", v)
}

func init() {
	rootCmd.AddCommand(newVersionCmd(os.Stdout))
}
