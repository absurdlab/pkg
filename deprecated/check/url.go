package check

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// URL is the entrypoint for checking a URL.
func URL(url *url.URL) *URLCheck {
	return (&URLCheck{value: url}).notNil()
}

// URLString is a convenience wrapper around URL.
func URLString(urlStr string) *URLCheck {
	u, err := url.Parse(urlStr)
	return &URLCheck{
		value: u,
		err:   err,
	}
}

type URLCheck struct {
	value *url.URL
	err   error
}

// HasScheme checks the url has one of the given schemes.
func (s *URLCheck) HasScheme(schemes ...string) *URLCheck {
	if s.shouldSkip() {
		return s
	}

	scheme := strings.ToLower(s.value.Scheme)
	for _, each := range schemes {
		if strings.ToLower(each) == scheme {
			return s
		}
	}

	s.err = fmt.Errorf("url scheme does not match one of [%s]", strings.Join(schemes, ", "))

	return s
}

// NotHaveScheme checks the url does not have one of the given schemes.
func (s *URLCheck) NotHaveScheme(schemes ...string) *URLCheck {
	if s.shouldSkip() {
		return s
	}

	scheme := strings.ToLower(s.value.Scheme)
	for _, each := range schemes {
		if strings.ToLower(each) == scheme {
			s.err = fmt.Errorf("url should not have %s scheme, but does", each)
			return s
		}
	}

	return s
}

// HasHost checks if the url hostname matches any of the given hosts.
func (s *URLCheck) HasHost(hosts ...string) *URLCheck {
	if s.shouldSkip() {
		return s
	}

	host := strings.ToLower(s.value.Hostname())
	for _, each := range hosts {
		if strings.ToLower(each) == host {
			return s
		}
	}

	s.err = fmt.Errorf("url hostname does not match one of [%s]", strings.Join(hosts, ", "))

	return s
}

// NotHaveHost checks if the url hostname does not match any of the given hosts.
func (s *URLCheck) NotHaveHost(hosts ...string) *URLCheck {
	if s.shouldSkip() {
		return s
	}

	host := strings.ToLower(s.value.Hostname())
	for _, each := range hosts {
		if strings.ToLower(each) == host {
			s.err = fmt.Errorf("url should not have %s as hostname, but does", each)
			return s
		}
	}

	return s
}

// IsLocalhost checks the hostname counts as localhost
func (s *URLCheck) IsLocalhost() *URLCheck {
	return s.HasHost("127.0.0.1", "localhost")
}

// IsNotLocalhost checks the hostname does not count as localhost
func (s *URLCheck) IsNotLocalhost() *URLCheck {
	return s.NotHaveHost("127.0.0.1", "localhost")
}

// NotHaveFragment checks the url does not have a fragment component.
func (s *URLCheck) NotHaveFragment() *URLCheck {
	if s.shouldSkip() {
		return s
	}

	if len(s.value.RawFragment) > 0 || len(s.value.Fragment) > 0 {
		s.err = errors.New("url should not have fragment")
	}

	return s
}

// NotHaveQuery checks the url does not have a query component.
func (s *URLCheck) NotHaveQuery() *URLCheck {
	if s.shouldSkip() {
		return s
	}

	if len(s.value.RawQuery) > 0 || len(s.value.Query()) > 0 {
		s.err = errors.New("url should not have query")
	}

	return s
}

// That performs a custom check on the url.
func (s *URLCheck) That(customCheck func(u *url.URL) error) *URLCheck {
	if s.shouldSkip() {
		return s
	}

	s.err = customCheck(s.value)

	return s
}

// Error terminates the check and returns the error.
func (s *URLCheck) Error() error {
	return s.err
}

func (s *URLCheck) shouldSkip() bool {
	return s.err != nil
}

func (s *URLCheck) notNil() *URLCheck {
	if s == nil {
		s.err = errors.New("url is nil")
	}
	return s
}
