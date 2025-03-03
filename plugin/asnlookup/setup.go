package asnlookup

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	var dbPath string

	for c.Next() {
		if !c.NextArg() {
			return plugin.Error(pluginName, c.ArgErr())
		}
		if dbPath != "" {
			return plugin.Error(pluginName, c.Errf("multiple databases not supported"))
		}
		dbPath = c.Val()
		if len(c.RemainingArgs()) != 0 {
			return plugin.Error(pluginName, c.ArgErr())
		}
	}

	asnLookup, err := NewASNLookup(dbPath)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		asnLookup.Next = next
		return asnLookup
	})

	return nil
}
