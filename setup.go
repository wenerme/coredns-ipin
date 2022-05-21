package ipin

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	caddy.RegisterPlugin(Name, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	ipin := IpInName{}

	c.Next() // ipin
	for c.NextBlock() {
		x := c.Val()
		switch x {
		case "fallback":
			ipin.Fallback = true
		case "ttl":
			args := c.RemainingArgs()
			if len(args) < 1 {
				return c.Errf("ttl needs a time in second")
			}
			ttl, err := strconv.Atoi(args[0])
			if err != nil {
				return c.Errf("ttl provided is invalid")
			}
			ipin.Ttl = uint32(ttl)
		default:
			return plugin.Error(Name, c.Errf("unexpected '%v' command", x))
		}
	}
	if c.NextArg() {
		return plugin.Error(Name, c.ArgErr())
	}

	dnsserver.
		GetConfig(c).
		AddPlugin(func(next plugin.Handler) plugin.Handler {
			ipin.Next = next
			return ipin
		})

	return nil
}
