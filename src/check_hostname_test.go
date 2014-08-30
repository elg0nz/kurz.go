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
  t.Logf("Will short URLs from following sites: %q. Length %d", allowedHostNames, len(allowedHostNames))
  var checkUrl, _ = url.Parse("http://www.somewhere.com")

  if check_hostname(allowedHostNames, checkUrl) != true {
    t.Error("Any URL should pass with empty hostname list")
  }
}

func TestCheckHostnameWithUrlWithAllowedHostname(t *testing.T) {
  logBackend := logging.NewChannelMemoryBackend(999)
  logging.SetBackend(logBackend)
  var allowedHostNames = []string{"www.somewhere.com"}
  t.Logf("Will short URLs from following sites: %q. Length %d", allowedHostNames, len(allowedHostNames))
  var checkUrl, _ = url.Parse("http://www.somewhere.com/example_location")

  if check_hostname(allowedHostNames, checkUrl) != true {
    t.Error("Any URL should pass with hostname in accepted list")
  }
}

func TestCheckHostnameWithUrlWithAllowedHostnames(t *testing.T) {
  logBackend := logging.NewChannelMemoryBackend(999)
  logging.SetBackend(logBackend)
  var allowedHostNames = []string{"www.somewhere.com", "www.elsewhere.com"}
  t.Logf("Will short URLs from following sites: %q. Length %d", allowedHostNames, len(allowedHostNames))
  var checkUrl, _ = url.Parse("http://www.somewhere.com/example_location")

  if check_hostname(allowedHostNames, checkUrl) != true {
    t.Error("Any URL should pass with hostname in accepted list")
  }
}

func TestCheckHostnameWithUrlWithoutAllowedHostnames(t *testing.T) {
  logBackend := logging.NewChannelMemoryBackend(999)
  logging.SetBackend(logBackend)
  var allowedHostNames = []string{"www.somewhere.com", "www.elsewhere.com"}
  t.Logf("Will short URLs from following sites: %q. Length %d", allowedHostNames, len(allowedHostNames))
  var checkUrl, _ = url.Parse("http://www.nowhere.com/example_location")

  if check_hostname(allowedHostNames, checkUrl) != false {
    t.Error("Any URL should not pass without hostname in accepted list")
  }
}

func TestCheckHostnameWithHttpsNotExistingOnListWithoutHttps(t *testing.T) {
  logBackend := logging.NewChannelMemoryBackend(999)
  logging.SetBackend(logBackend)
  var allowedHostNames = []string{"www.somewhere.com", "www.elsewhere.com"}
  t.Logf("Will short URLs from following sites: %q. Length %d", allowedHostNames, len(allowedHostNames))
  var checkUrl, _ = url.Parse("https://www.nowhere.com/example_location")

  if check_hostname(allowedHostNames, checkUrl) != false {
    t.Error("Any URL should not pass without hostname in accepted list")
  }
}

func TestCheckHostnameWithHttpsOnListWithoutHttps(t *testing.T) {
  logBackend := logging.NewChannelMemoryBackend(999)
  logging.SetBackend(logBackend)
  var allowedHostNames = []string{"www.somewhere.com", "www.elsewhere.com"}
  t.Logf("Will short URLs from following sites: %q. Length %d", allowedHostNames, len(allowedHostNames))
  var checkUrl, _ = url.Parse("https://www.somewhere.com/example_location")

  if check_hostname(allowedHostNames, checkUrl) != true {
    t.Error("Any URL should pass with hostname in accepted list")
  }
}
