package torque

import (
	"net"
	"net/url"
)

// HostPortFlag handles converting a "host:port" string into a full-fledged
// url.URL
type HostPortFlag url.URL

// String calls url.URL's String() method
func (hpf *HostPortFlag) String() string {
	u := url.URL(*hpf)
	return u.String()
}

// Set parses the host:port string into a valid url.URL
func (hpf *HostPortFlag) Set(value string) error {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return err
	}
	hp := net.JoinHostPort(host, port)
	*hpf = HostPortFlag(url.URL{
		Scheme: Scheme,
		Host:   hp,
	})
	return nil
}
