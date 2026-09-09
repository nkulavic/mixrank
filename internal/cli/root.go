package cli

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nkulavic/mixrank"
	"github.com/nkulavic/mixrank/catalog"
	"github.com/nkulavic/mixrank/internal/credentials"
	"github.com/nkulavic/mixrank/mcpserver"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

var Version = "0.5.0"

func New() *cobra.Command {
	root := &cobra.Command{Use: "mixrank", Short: "MixRank API, research skills and MCP toolkit", Version: Version, SilenceUsage: true, SilenceErrors: true}
	root.AddCommand(auth(), endpoints(), api(), request(), mcpCommand(), doctor(), setupCommand(), personalizeCommand(), profilesCommand(), downloadCommand(), validationCommand(), contactsCommand())
	return root
}
func client(ctx context.Context) (*mixrank.Client, error) {
	key, _, e := credentials.Resolve()
	if e != nil {
		return nil, e
	}
	return mixrank.New(key, mixrank.Options{Retries: 2})
}
func jsonOut(cmd *cobra.Command, v any) error {
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
func endpoints() *cobra.Command {
	var index string
	c := &cobra.Command{Use: "endpoints [operation]", Short: "Inspect all operations or one complete contract", Args: cobra.MaximumNArgs(1), RunE: func(c *cobra.Command, a []string) error {
		if index != "" {
			for _, v := range []string{"person2", "companies", "jobs", "org_name", "industries"} {
				if index == v {
					b, e := catalog.Files.ReadFile("mappings/" + index + ".json")
					if e != nil {
						return e
					}
					_, e = c.OutOrStdout().Write(b)
					return e
				}
			}
			return errors.New("unknown index")
		}
		if len(a) > 0 {
			op, e := catalog.Lookup(a[0])
			if e != nil {
				return e
			}
			return jsonOut(c, op)
		}
		return jsonOut(c, catalog.Read())
	}}
	c.Flags().StringVar(&index, "mapping", "", "Elasticsearch index mapping")
	return c
}
func auth() *cobra.Command {
	p := &cobra.Command{Use: "auth", Short: "Manage toolkit provider credentials"}
	var stdin bool
	login := &cobra.Command{Use: "login", Short: "Save a masked API key in the OS credential store", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		var b []byte
		var e error
		if stdin {
			b, e = io.ReadAll(io.LimitReader(c.InOrStdin(), 4097))
			if len(b) > 4096 {
				return errors.New("credential input too long")
			}
		} else {
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return errors.New("no interactive terminal; pipe the key with --stdin")
			}
			fmt.Fprint(c.ErrOrStderr(), "MixRank API key: ")
			b, e = term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Fprintln(c.ErrOrStderr())
		}
		if e != nil {
			return errors.New("credential input failed")
		}
		defer func() {
			for i := range b {
				b[i] = 0
			}
		}()
		if e = credentials.Save(strings.TrimSpace(string(b))); e != nil {
			return e
		}
		return jsonOut(c, map[string]any{"saved": true, "storage": "os-vault", "service": credentials.Service, "account": credentials.Account})
	}}
	login.Flags().BoolVar(&stdin, "stdin", false, "Read key from stdin; never pass it as an argument")
	status := &cobra.Command{Use: "status", Short: "Report credential availability without revealing it", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		_, source, e := credentials.Resolve()
		return jsonOut(c, map[string]any{"configured": e == nil, "source": source, "service": credentials.Service, "account": credentials.Account})
	}}
	logout := &cobra.Command{Use: "logout", Short: "Delete MixRank's own OS-vault entry", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if e := credentials.Delete(); e != nil {
			return e
		}
		return jsonOut(c, map[string]any{"deleted": true, "environment_override_present": os.Getenv("MIXRANK_API_KEY") != ""})
	}}
	p.AddCommand(login, status, logout)
	places := &cobra.Command{Use: "google-places", Short: "Manage Google Places fallback credentials"}
	var placesStdin bool
	placesLogin := &cobra.Command{Use: "login", Short: "Save a masked Google Places API key in the OS credential store", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		var b []byte
		var e error
		if placesStdin {
			b, e = io.ReadAll(io.LimitReader(c.InOrStdin(), 4097))
			if len(b) > 4096 {
				return errors.New("credential input too long")
			}
		} else {
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return errors.New("no interactive terminal; pipe the key with --stdin")
			}
			fmt.Fprint(c.ErrOrStderr(), "Google Places API key: ")
			b, e = term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Fprintln(c.ErrOrStderr())
		}
		if e != nil {
			return errors.New("credential input failed")
		}
		defer func() {
			for i := range b {
				b[i] = 0
			}
		}()
		if e = credentials.SaveGooglePlaces(strings.TrimSpace(string(b))); e != nil {
			return e
		}
		return jsonOut(c, map[string]any{"saved": true, "storage": "os-vault", "service": credentials.PlacesService, "account": credentials.PlacesAccount})
	}}
	placesLogin.Flags().BoolVar(&placesStdin, "stdin", false, "Read key from stdin; never pass it as an argument")
	placesStatus := &cobra.Command{Use: "status", Short: "Report Google Places credential availability without revealing it", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		_, source, e := credentials.ResolveGooglePlaces()
		return jsonOut(c, map[string]any{"configured": e == nil, "source": source, "service": credentials.PlacesService, "account": credentials.PlacesAccount})
	}}
	placesLogout := &cobra.Command{Use: "logout", Short: "Delete the toolkit's Google Places OS-vault entry", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if e := credentials.DeleteGooglePlaces(); e != nil {
			return e
		}
		return jsonOut(c, map[string]any{"deleted": true, "environment_override_present": os.Getenv("GOOGLE_PLACES_API_KEY") != "" || os.Getenv("GOOGLE_MAPS_API_KEY") != ""})
	}}
	places.AddCommand(placesLogin, placesStatus, placesLogout)
	p.AddCommand(places)
	return p
}

type requestFlags struct {
	query                        []string
	body, upload, output, format string
	refresh                      bool
	pages                        int
	maxBytes                     int64
}

func (f *requestFlags) flags(c *cobra.Command) {
	c.Flags().StringArrayVarP(&f.query, "query", "q", nil, "Query KEY=VALUE (repeatable)")
	c.Flags().StringVar(&f.body, "body", "", "JSON file path, or - for stdin")
	c.Flags().StringVar(&f.upload, "upload", "", "UTF-8 bulk upload file path, or - for stdin")
	c.Flags().StringVarP(&f.output, "output", "o", "", "Stream raw response to a file (private permissions)")
	c.Flags().StringVar(&f.format, "format", "json", "json, jsonl, csv, raw")
	c.Flags().BoolVar(&f.refresh, "refresh", false, "Explicitly allow requested live-fetch strategy")
	c.Flags().IntVar(&f.pages, "max-pages", 1, "Bounded offset pagination, 1..1000")
	c.Flags().Int64Var(&f.maxBytes, "max-bytes", mixrank.DefaultMaxResponseBytes, "Maximum JSON response bytes")
}
func (f *requestFlags) input(c *cobra.Command) (mixrank.Request, func(), error) {
	in := mixrank.Request{Query: url.Values{}, Refresh: f.refresh}
	closers := []io.Closer{}
	cleanup := func() {
		for _, r := range closers {
			r.Close()
		}
	}
	for _, kv := range f.query {
		k, v, ok := strings.Cut(kv, "=")
		if !ok || k == "" {
			return in, cleanup, errors.New("query requires KEY=VALUE")
		}
		in.Query.Add(k, v)
	}
	if f.body == "-" && f.upload == "-" {
		return in, cleanup, errors.New("stdin cannot supply both body and upload")
	}
	read := func(p string) (io.Reader, error) {
		if p == "-" {
			return c.InOrStdin(), nil
		}
		r, e := os.Open(p)
		if e == nil {
			closers = append(closers, r)
		}
		return r, e
	}
	if f.body != "" {
		r, e := read(f.body)
		if e != nil {
			return in, cleanup, e
		}
		b, e := io.ReadAll(io.LimitReader(r, 16<<20+1))
		if e != nil || len(b) > 16<<20 {
			return in, cleanup, errors.New("JSON input exceeds 16 MiB or cannot be read")
		}
		if !json.Valid(b) {
			return in, cleanup, errors.New("invalid JSON body")
		}
		in.Body = b
	}
	if f.upload != "" {
		r, e := read(f.upload)
		if e != nil {
			return in, cleanup, e
		}
		in.Upload = r
		in.Filename = f.upload
	}
	return in, cleanup, nil
}
func api() *cobra.Command {
	a := &cobra.Command{Use: "api", Short: "Named commands generated from every catalog operation"}
	for _, op := range catalog.Read().Operations {
		f := &requestFlags{}
		values := map[string]*string{}
		c := &cobra.Command{Use: strings.ReplaceAll(op.ID, "_", "-"), Short: op.Summary, Long: op.Summary + "\n\n" + op.Notes, Args: cobra.NoArgs}
		f.flags(c)
		for _, p := range op.Parameters {
			if p.In == "file" {
				continue
			}
			name := strings.ReplaceAll(p.Name, "_", "-")
			if c.Flags().Lookup(name) != nil {
				continue
			}
			values[p.Name] = c.Flags().String(name, "", p.Description)
		}
		c.RunE = func(c *cobra.Command, _ []string) error {
			in, close, e := f.input(c)
			defer close()
			if e != nil {
				return e
			}
			in.Parameters = map[string]string{}
			for k, v := range values {
				if c.Flags().Changed(strings.ReplaceAll(k, "_", "-")) {
					in.Parameters[k] = *v
				}
			}
			cl, e := client(c.Context())
			if e != nil {
				return e
			}
			if f.pages > 1 {
				if f.output != "" || f.format == "csv" || f.format == "raw" {
					return errors.New("multi-page mode emits JSONL pages; use one-page files for binary/CSV")
				}
				return cl.Pages(c.Context(), op.ID, in, f.pages, func(page int, v any) error {
					return json.NewEncoder(c.OutOrStdout()).Encode(map[string]any{"page": page, "data": v})
				})
			}
			if f.pages != 1 {
				return errors.New("max-pages must be positive")
			}
			r, e := cl.Call(c.Context(), op.ID, in)
			if e != nil {
				return e
			}
			defer r.Body.Close()
			return output(c, r, f)
		}
		a.AddCommand(c)
	}
	return a
}
func request() *cobra.Command {
	f := &requestFlags{}
	c := &cobra.Command{Use: "request METHOD /path", Short: "Raw API access with shared credential and concurrency safeguards", Args: cobra.ExactArgs(2), RunE: func(c *cobra.Command, a []string) error {
		in, close, e := f.input(c)
		defer close()
		if e != nil {
			return e
		}
		if f.upload != "" {
			return errors.New("use a named bulk operation for multipart uploads")
		}
		cl, e := client(c.Context())
		if e != nil {
			return e
		}
		var body io.Reader
		typ := ""
		if in.Body != nil {
			body = strings.NewReader(string(in.Body))
			typ = "application/json"
		}
		r, e := cl.Raw(c.Context(), a[0], a[1], in.Query, body, typ, in.Refresh)
		if e != nil {
			return e
		}
		defer r.Body.Close()
		return output(c, r, f)
	}}
	f.flags(c)
	return c
}
func output(c *cobra.Command, r *mixrank.Response, f *requestFlags) error {
	if f.output != "" {
		file, e := os.OpenFile(f.output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		_, e = io.Copy(file, r.Body)
		ce := file.Close()
		if e != nil {
			return errors.New("download interrupted; partial file retained")
		}
		return ce
	}
	if f.format == "raw" {
		_, e := io.Copy(c.OutOrStdout(), r.Body)
		return e
	}
	v, e := r.JSON(f.maxBytes)
	if e != nil {
		return e
	}
	return export(c.OutOrStdout(), v, f.format)
}
func export(w io.Writer, v any, format string) error {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	case "jsonl", "csv":
	default:
		return errors.New("format must be json, jsonl, csv, or raw")
	}
	rows, ok := v.([]any)
	if !ok {
		if m, yes := v.(map[string]any); yes {
			rows, ok = m["results"].([]any)
		}
		if !ok {
			rows = []any{v}
		}
	}
	if format == "jsonl" {
		for _, r := range rows {
			if e := json.NewEncoder(w).Encode(r); e != nil {
				return e
			}
		}
		return nil
	}
	columns := map[string]bool{}
	for _, r := range rows {
		m, ok := r.(map[string]any)
		if !ok {
			return errors.New("CSV requires object rows")
		}
		for k := range m {
			columns[k] = true
		}
	}
	keys := []string{}
	for k := range columns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	cw := csv.NewWriter(w)
	if e := cw.Write(keys); e != nil {
		return e
	}
	for _, r := range rows {
		m := r.(map[string]any)
		cells := []string{}
		for _, k := range keys {
			value := m[k]
			s, ok := value.(string)
			if !ok && value != nil {
				b, _ := json.Marshal(value)
				s = string(b)
			}
			if strings.HasPrefix(s, "=") || strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") || strings.HasPrefix(s, "@") || strings.HasPrefix(s, "\t") || strings.HasPrefix(s, "\r") {
				s = "'" + s
			}
			cells = append(cells, s)
		}
		if e := cw.Write(cells); e != nil {
			return e
		}
	}
	cw.Flush()
	return cw.Error()
}
func mcpCommand() *cobra.Command {
	var transport, addr string
	var origins []string
	c := &cobra.Command{Use: "mcp", Short: "Run stdio or authenticated Streamable HTTP MCP", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		switch transport {
		case "stdio":
			return mcpserver.StdioWithPlaces(c.Context(), client, placesClient)
		case "http":
			return mcpserver.ServeWithPlaces(c.Context(), addr, httpClientFactory, httpPlacesClient, mcpserver.HTTPOptions{BearerToken: os.Getenv("MIXRANK_MCP_TOKEN"), AllowedOrigins: origins})
		default:
			return errors.New("transport must be stdio or http")
		}
	}}
	c.Flags().StringVar(&transport, "transport", "stdio", "stdio or http")
	c.Flags().StringVar(&addr, "listen", "127.0.0.1:8080", "HTTP listen address (put TLS in front for remote access)")
	c.Flags().StringArrayVar(&origins, "origin", nil, "Allowed exact origin, repeatable")
	return c
}
func doctor() *cobra.Command {
	var live bool
	c := &cobra.Command{Use: "doctor", Short: "Local diagnostics; optional single /echo authentication check", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		_, source, e := credentials.Resolve()
		_, placesSource, placesErr := credentials.ResolveGooglePlaces()
		checks := map[string]any{"version": Version, "catalog_operations": len(catalog.Read().Operations), "credential_available": e == nil, "credential_source": source, "places_credential_available": placesErr == nil, "places_credential_source": placesSource}
		if live {
			cl, e := client(c.Context())
			if e != nil {
				return e
			}
			ctx, cancel := context.WithTimeout(c.Context(), 20*time.Second)
			defer cancel()
			r, e := cl.GetEcho(ctx, mixrank.Request{})
			if e != nil {
				return e
			}
			defer r.Body.Close()
			v, e := r.JSON(4096)
			if e != nil {
				return e
			}
			checks["echo"] = v
		}
		return jsonOut(c, checks)
	}}
	c.Flags().BoolVar(&live, "live", false, "Make one authenticated /echo call")
	return c
}
func prompt(c *cobra.Command, label string) (string, error) {
	fmt.Fprint(c.ErrOrStderr(), label)
	s, e := bufio.NewReader(c.InOrStdin()).ReadString('\n')
	return strings.TrimSpace(s), e
}

// ExitCode distinguishes invalid input, credentials, HTTP rejection and uncertain outcomes.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var api *mixrank.APIError
	if errors.As(err, &api) {
		if api.Uncertain {
			return 6
		}
		if api.Status == 401 || api.Status == 403 {
			return 3
		}
		if api.Status == 429 {
			return 4
		}
		return 5
	}
	if errors.Is(err, context.Canceled) {
		return 130
	}
	if strings.Contains(err.Error(), "credential") {
		return 3
	}
	return 2
}

var _ = http.StatusOK

func httpClientFactory(ctx context.Context) (*mixrank.Client, error) {
	key := os.Getenv("MIXRANK_API_KEY")
	if key == "" {
		return nil, errors.New("HTTP MCP requires upstream MIXRANK_API_KEY from server environment")
	}
	return mixrank.New(key, mixrank.Options{Retries: 2, InMemoryCoordination: true})
}

func placesClient(ctx context.Context) (*mixrank.PlacesClient, error) {
	key, _, err := credentials.ResolveGooglePlaces()
	if err != nil {
		return nil, err
	}
	return mixrank.NewPlacesClient(key, mixrank.PlacesOptions{})
}

func httpPlacesClient(ctx context.Context) (*mixrank.PlacesClient, error) {
	key := os.Getenv("GOOGLE_PLACES_API_KEY")
	if key == "" {
		key = os.Getenv("GOOGLE_MAPS_API_KEY")
	}
	if key == "" {
		return nil, errors.New("HTTP MCP places fallback requires GOOGLE_PLACES_API_KEY from server environment")
	}
	return mixrank.NewPlacesClient(key, mixrank.PlacesOptions{})
}
