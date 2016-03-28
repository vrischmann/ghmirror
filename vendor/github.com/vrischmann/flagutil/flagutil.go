// Package flagutil provides a collection of types implementing the flag.Value interface.
package flagutil

import (
	"net"
	"net/url"
	"strings"
)

// NetworkAddresses is a slice of string that have been validated as valid network addresses.
// Use it as a flag value when you want to pass a comma-separated list of strings to a flag
// and have it to be automatically parsed into a slice and validated as valid network addresses.
type NetworkAddresses []string

// Set implements the flag.Value interface. It parses the string as a comma-separated string.
// Additionally, each value is then passed to net.SplitHostPort for validation that it's a correct network address.
func (a *NetworkAddresses) Set(s string) error {
	for _, t := range strings.Split(s, ",") {
		_, _, err := net.SplitHostPort(t)
		if err != nil {
			return err
		}

		*a = append(*a, t)
	}

	return nil
}

// String implements the flag.Value interface. It returns the slice as a comma-separated string.
func (a NetworkAddresses) String() string {
	return strings.Join(a, ",")
}

// StringSlice returns the network addresses as a slice of string.
func (a NetworkAddresses) StringSlice() []string {
	var res []string
	for _, el := range a {
		res = append(res, el)
	}
	return res
}

// Strings is a slice of string.
// Use it as a flag value when you want to pass a comma-separated list of strings to a flag
// and have it to be automatically parsed into a slice.
type Strings []string

// Set implements the flag.Value interface. It parses the string as a comma-separated string.
func (s *Strings) Set(str string) error {
	for _, t := range strings.Split(str, ",") {
		*s = append(*s, t)
	}

	return nil
}

// String implements the flag.Value interface. It returns the slice as a comma-separated string.
func (s Strings) String() string {
	return strings.Join(s, ",")
}

// URL is a wrapper around a url.URL type.
// Use it as a flag value when you want to validate a flag as a valid URL.
type URL struct {
	url.URL
	parsed bool
}

// Set implements the flag.Value interface. It parses the string as a url.URL value.
func (u *URL) Set(s string) error {
	url, err := url.Parse(s)
	if err != nil {
		return err
	}

	u.URL = *url
	u.parsed = true

	return nil
}

// String implements the flag.Value interface. It returns the underlying url.URL as a string.
func (u URL) String() string {
	return u.URL.String()
}

// IsValid returns true if the URL object is valid, false otherwise. Valid here means it has correctly parsed an URL.
func (u URL) IsValid() bool {
	return u.parsed
}

// URLs is a slice of url.URL.
// Use it as a flag value when you want to pass a comma-separated list of strings to a flag
// and have it to be automatically parsed into a slice and validated as valid url.URL.
type URLs []url.URL

// Set implements the flag.Value interface. It parses the string as a comma-separated string.
func (u *URLs) Set(s string) error {
	for _, t := range strings.Split(s, ",") {
		url, err := url.Parse(t)
		if err != nil {
			return err
		}

		*u = append(*u, *url)
	}

	return nil
}

// String implements the flag.Value interface. It returns the string representation of each url.URL as a comma-separated string.
func (u URLs) String() string {
	var r []string

	for _, url := range u {
		r = append(r, url.String())
	}

	return strings.Join(r, ",")
}
