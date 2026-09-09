package mixrank

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nkulavic/mixrank/catalog"
	"net/url"
	"strconv"
)

// Pages uses only documented offset/page_size pagination. maxPages is mandatory;
// each response is delivered separately, retaining partial progress on failure.
func (c *Client) Pages(ctx context.Context, id string, in Request, maxPages int, visit func(int, any) error) error {
	op, e := catalog.Lookup(id)
	if e != nil {
		return e
	}
	if op.Pagination != "offset" {
		return errors.New("operation does not document offset pagination; use explicit DSL or provider cursor")
	}
	if maxPages < 1 || maxPages > 1000 {
		return errors.New("max pages must be 1..1000")
	}
	q := url.Values{}
	for k, v := range in.Query {
		q[k] = append([]string(nil), v...)
	}
	p := map[string]string{}
	for k, v := range in.Parameters {
		p[k] = v
	}
	in.Parameters = p
	in.Query = q
	size := 100
	offset := 0
	for _, par := range op.Parameters {
		if par.Name == "page_size" {
			if n, e := strconv.Atoi(par.Default); e == nil {
				size = n
			}
		}
	}
	for k, dst := range map[string]*int{"page_size": &size, "offset": &offset} {
		s := p[k]
		if s == "" {
			s = q.Get(k)
		}
		if s != "" {
			n, e := strconv.Atoi(s)
			if e != nil || n < 0 {
				return fmt.Errorf("invalid %s", k)
			}
			*dst = n
		}
	}
	if size < 1 {
		return errors.New("page_size must be positive")
	}
	for page := 0; page < maxPages; page++ {
		p["offset"] = strconv.Itoa(offset)
		p["page_size"] = strconv.Itoa(size)
		r, e := c.Call(ctx, id, in)
		if e != nil {
			return e
		}
		v, e := r.JSON(DefaultMaxResponseBytes)
		r.Body.Close()
		if e != nil {
			return e
		}
		if e = visit(page, v); e != nil {
			return e
		}
		n := -1
		switch x := v.(type) {
		case []any:
			n = len(x)
		case map[string]any:
			for _, key := range []string{"results", "data"} {
				if a, ok := x[key].([]any); ok {
					n = len(a)
					break
				}
			}
			if total, ok := x["total"].(json.Number); ok {
				if t, e := total.Int64(); e == nil && int64(offset+size) >= t {
					return nil
				}
			}
		}
		if n < 0 {
			return errors.New("unknown pagination response shape; explicit next request required")
		}
		if n < size {
			return nil
		}
		offset += size
	}
	return nil
}
