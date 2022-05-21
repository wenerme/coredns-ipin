package ipin_test

import (
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
	"github.com/wenerme/coredns-ipin"
	"golang.org/x/net/context"
)

func TestWhoami(t *testing.T) {
	wh := ipin.IpInName{ Ttl: uint32(86400) }

	tests := []struct {
		qname         string
		qtype         uint16
		expectedCode  int
		expectedReply []string // ownernames for the records in the additional section.
		expectedTtl  uint32
		expectedErr   error
	}{
		{
			qname:         "192-168-1-2-80.example.org",
			qtype:         dns.TypeA,
			expectedCode:  dns.RcodeSuccess,
			expectedReply: []string{"192-168-1-2-80.example.org.", "_port.192-168-1-2-80.example.org."},
			expectedTtl:   uint32(86400),
			expectedErr:   nil,
		},
	}

	ctx := context.TODO()

	for i, tc := range tests {
		req := new(dns.Msg)
		req.SetQuestion(dns.Fqdn(tc.qname), tc.qtype)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		code, err := wh.ServeDNS(ctx, rec, req)

		if err != tc.expectedErr {
			t.Errorf("Test %d: Expected error %v, but got %v", i, tc.expectedErr, err)
		}
		if code != int(tc.expectedCode) {
			t.Errorf("Test %d: Expected status code %d, but got %d", i, tc.expectedCode, code)
		}
		if len(tc.expectedReply) != 0 {
			actual := rec.Msg.Answer[0].Header().Name
			expected := tc.expectedReply[0]
			if actual != expected {
				t.Errorf("Test %d: Expected answer %s, but got %s", i, expected, actual)
			}

			actual = rec.Msg.Extra[0].Header().Name
			expected = tc.expectedReply[1]
			if actual != expected {
				t.Errorf("Test %d: Expected answer %s, but got %s", i, expected, actual)
			}

			if rec.Msg.Extra[0].Header().Ttl != tc.expectedTtl {
				t.Errorf("Test %d: Expected answer %d, but got %d", i, tc.expectedTtl, rec.Msg.Extra[0].Header().Ttl)
			}
		}
	}
}
