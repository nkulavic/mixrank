package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/nkulavic/mixrank"
	"github.com/spf13/cobra"
)

func validationFlags(c *cobra.Command, o *mixrank.ContactValidationOptions) {
	c.Flags().StringVar(&o.Strategy, "validation-strategy", "besteffort", "Email validation strategy: cached, besteffort, strict, fetch")
	c.Flags().StringVar(&o.MaxAge, "validation-maxage", "", "Maximum cached validation age in seconds")
	c.Flags().DurationVar(&o.Wait, "validation-wait", 30*time.Second, "Wait for bulk validation (0 submits/checks once; up to 10m)")
}
func newContactValidationCommand(factory func(context.Context) (*mixrank.Client, error)) *cobra.Command {
	var input, output string
	var opts mixrank.ContactValidationOptions
	var timeout time.Duration
	c := &cobra.Command{Use: "validate", Short: "Validate emails in a contact report, or resume its existing bulk job", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if input == "" {
			return errors.New("--input contact report is required")
		}
		if timeout <= 0 || timeout > 30*time.Minute {
			return errors.New("timeout must be greater than zero and at most 30m")
		}
		var reader io.Reader = c.InOrStdin()
		if input != "-" {
			f, e := os.Open(input)
			if e != nil {
				return e
			}
			defer f.Close()
			reader = f
		}
		b, e := io.ReadAll(io.LimitReader(reader, 4<<20+1))
		if e != nil || len(b) > 4<<20 {
			return errors.New("report must fit within 4 MiB")
		}
		var report mixrank.ContactReport
		if json.Unmarshal(b, &report) != nil || report.Companies == nil {
			return errors.New("expected a contact report JSON object")
		}
		writer := c.OutOrStdout()
		written := false
		if output != "" {
			f, e := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
			if e != nil {
				return e
			}
			defer func() {
				f.Close()
				if !written {
					os.Remove(output)
				}
			}()
			writer = f
		}
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()
		cl, e := factory(ctx)
		if e != nil {
			return e
		}
		err := cl.ValidateContacts(ctx, &report, opts)
		enc := json.NewEncoder(writer)
		enc.SetIndent("", "  ")
		if e := enc.Encode(report); e != nil {
			return e
		}
		written = true
		return err
	}}
	c.Flags().StringVar(&input, "input", "", "Existing contact report file, or - for stdin")
	c.Flags().StringVarP(&output, "output", "o", "", "Write private JSON; refuses existing files")
	c.Flags().IntVar(&opts.MaxRequests, "max-requests", 30, "Maximum validation submit/status/download requests")
	c.Flags().DurationVar(&timeout, "timeout", 10*time.Minute, "Overall timeout")
	validationFlags(c, &opts)
	return c
}
