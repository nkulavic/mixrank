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

func contactsCommand() *cobra.Command { return newContactsCommand(client) }

func newContactsCommand(factory func(context.Context) (*mixrank.Client, error)) *cobra.Command {
	var opts mixrank.ContactOptions
	var input, output, domain, name string
	var ids []string
	var timeout time.Duration
	c := &cobra.Command{Use: "contacts", Short: "Find current company contacts and retrieve available business emails and direct dials", Long: "Discover current owners and relevant managers, append B2B emails and direct dials, and return one JSON report. Accepts a company array, a companies export, or raw Elasticsearch company hits. Keeps missing contacts, inferred employment, stale profiles, and partial failures explicit. Use --validate-emails for bulk deliverability evidence. Does not send messages.", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if timeout <= 0 || timeout > 30*time.Minute {
			return errors.New("timeout must be greater than zero and at most 30m")
		}
		if input != "" && (len(ids) > 0 || domain != "" || name != "") {
			return errors.New("use --companies or a single --company-id/--domain target")
		}
		if input != "" {
			var r io.Reader = c.InOrStdin()
			if input != "-" {
				f, e := os.Open(input)
				if e != nil {
					return e
				}
				defer f.Close()
				r = f
			}
			companies, e := mixrank.ParseCompanyTargets(r)
			if e != nil {
				return e
			}
			opts.Companies = companies
		} else {
			opts.Companies = []mixrank.CompanyTarget{{CompanyIDs: ids, Domain: domain, Name: name}}
		}
		writer := c.OutOrStdout()
		var file *os.File
		written := false
		if output != "" {
			var e error
			file, e = os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			if e != nil {
				return e
			}
			defer func() {
				file.Close()
				if !written {
					os.Remove(output)
				}
			}()
			writer = file
		}
		cl, e := factory(c.Context())
		if e != nil {
			return e
		}
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()
		report, workflowErr := cl.CompanyContacts(ctx, opts)
		if report != nil {
			enc := json.NewEncoder(writer)
			enc.SetIndent("", "  ")
			if e = enc.Encode(report); e != nil {
				return e
			}
			written = true
		}
		return workflowErr
	}}
	c.Flags().IntVar(&opts.Concurrency, "concurrency", 4, "Maximum simultaneous company/person requests (1..16)")
	c.Flags().StringVar(&opts.ContactFilter, "contact-filter", "all", "Return all, any, email, phone, both, none, or valid-email contacts")
	c.Flags().BoolVar(&opts.ValidateEmails, "validate-emails", false, "Validate returned work emails through a deduplicated bulk job")
	validationFlags(c, &opts.ValidationOptions)
	c.AddCommand(newContactValidationCommand(factory))
	c.Flags().StringVar(&input, "companies", "", "Company JSON file, or - for stdin")
	c.Flags().StringArrayVar(&ids, "company-id", nil, "MixRank company ID (repeatable for aliases of one company)")
	c.Flags().StringVar(&domain, "domain", "", "Company website domain")
	c.Flags().StringVar(&name, "name", "", "Company display name; requires an ID or domain")
	c.Flags().StringSliceVar(&opts.Roles, "roles", nil, "Ordered role phrases; defaults to owners, executives and managers")
	c.Flags().IntVar(&opts.MaxContacts, "max-contacts", 2, "Stop enrichment after this many contacts have requested channels (1..5); all also includes other discovered people")
	c.Flags().IntVar(&opts.MaxCandidates, "max-candidates", 10, "Maximum candidates to consider per company (contacts..25)")
	c.Flags().IntVar(&opts.MaxRequests, "max-requests", 100, "Maximum upstream workflow calls (1..250); partial results retained")
	c.Flags().BoolVar(&opts.EmailsOnly, "emails-only", false, "Request B2B emails without direct dials")
	c.Flags().StringVarP(&output, "output", "o", "", "Write private JSON; refuses existing files")
	c.Flags().DurationVar(&timeout, "timeout", 10*time.Minute, "Overall workflow timeout")
	return c
}
