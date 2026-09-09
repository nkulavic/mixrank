package cli

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

func validationCommand() *cobra.Command {
	root := &cobra.Command{Use: "validation", Short: "Inspect or acknowledge asynchronous one-off validation"}
	root.AddCommand(&cobra.Command{Use: "status", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		cl, e := client(c.Context())
		if e != nil {
			return e
		}
		p, e := cl.ValidationStatus()
		if e != nil {
			return e
		}
		return jsonOut(c, p)
	}})
	var headers, body string
	ack := &cobra.Command{Use: "ack", Short: "Verify and acknowledge a matching signed completion webhook", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if headers == "" || body == "" {
			return errors.New("--headers and --body are required")
		}
		hb, e := os.ReadFile(headers)
		if e != nil {
			return e
		}
		var h http.Header
		if e = json.Unmarshal(hb, &h); e != nil {
			return e
		}
		b, e := os.ReadFile(body)
		if e != nil {
			return e
		}
		cl, e := client(c.Context())
		if e != nil {
			return e
		}
		if e = cl.AcknowledgeValidationWebhook(c.Context(), h, b, time.Now(), 5*time.Minute); e != nil {
			return e
		}
		return jsonOut(c, map[string]bool{"acknowledged": true})
	}}
	ack.Flags().StringVar(&headers, "headers", "", "JSON header map from callback request")
	ack.Flags().StringVar(&body, "body", "", "Unmodified callback body file")
	root.AddCommand(ack)
	var confirmed bool
	clear := &cobra.Command{Use: "clear", Short: "Clear uncertain one-off state after independently confirming provider completion", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		cl, e := client(c.Context())
		if e != nil {
			return e
		}
		if e = cl.ClearValidation(c.Context(), confirmed); e != nil {
			return e
		}
		return jsonOut(c, map[string]bool{"cleared": true})
	}}
	clear.Flags().BoolVar(&confirmed, "confirmed-complete", false, "Completion has been confirmed with the callback receiver or provider")
	root.AddCommand(clear)
	return root
}
