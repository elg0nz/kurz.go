package main

import "net/url"

func check_hostname(allowedHostNames []string, input *url.URL) bool {
  log.Debug("Call to check if hostname in %q", allowedHostNames)

  if input == nil {
    return false
  }

  if len(allowedHostNames) == 0 {
    return true
  }

  return false
}
