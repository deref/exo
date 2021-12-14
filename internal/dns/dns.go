package dns

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	answer := []dns.RR{}
	for _, q := range r.Question {
		if q.Qtype == dns.TypeA && strings.HasSuffix(q.Name, ".localhost.") {
			answer = append(answer, &dns.A{
				Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
				A:   net.IPv4(127, 0, 0, 1),
			})
		}
	}
	m.Answer = answer
	w.WriteMsg(m)
}

func StartDNSServer() error {
	server := &dns.Server{Addr: ":4453", Net: "udp"}
	dns.HandleFunc(".", handleRequest)
	return server.ListenAndServe()
}
