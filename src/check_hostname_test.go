package main

import (
  "testing"
  "github.com/op/go-logging"
  "net/url"
)

func TestCheckHostnameWithoutUrl(t *testing.T) {
  logBackend := logging.NewChannelMemoryBackend(999)
  logging.SetBackend(logBackend)

  if check_hostname(nil, nil) != false {
    t.Error("nil check should return false")
  }
}

func TestCheckHostnameWithUrlWithoutHostnames(t *testing.T) {
  logBackend := logging.NewChannelMemoryBackend(999)
  logging.SetBackend(logBackend)
  var allowedHostNames = []string{}
  t.Logf("Will short URLs from following sites: %q. Length %i", allowedHostNames, len(allowedHostNames))
  var checkUrl, _ = url.Parse("http://www.somewhere.com")

  if check_hostname(allowedHostNames, checkUrl) != true {
    t.Error("Any URL should pass with empty hostname list")
  }
}
