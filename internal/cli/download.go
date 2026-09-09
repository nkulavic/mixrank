package cli

import (
	"errors"
	"github.com/nkulavic/mixrank"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func downloadCommand() *cobra.Command {
	var urlFile, output string
	var limit int64
	c := &cobra.Command{Use: "download", Short: "Download a completed job's HTTPS result without forwarding the API key", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if urlFile == "" || output == "" {
			return errors.New("--url-file and --output are required")
		}
		b, e := os.ReadFile(urlFile)
		if e != nil {
			return e
		}
		if len(b) > 16384 {
			return errors.New("URL file too large")
		}
		f, e := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		defer f.Close()
		return mixrank.Download(c.Context(), strings.TrimSpace(string(b)), f, limit)
	}}
	c.Flags().StringVar(&urlFile, "url-file", "", "Private file containing the returned download_url")
	c.Flags().StringVar(&output, "output", "", "New destination file")
	c.Flags().Int64Var(&limit, "max-bytes", 256<<20, "Maximum downloaded bytes")
	return c
}
