package main

import "net/url"

func check_hostname(allowedHostNames []string, input *url.URL) bool {
  log.Debug("Call to check if %s in %q", input.Host, allowedHostNames)

  if input == nil {
    return false
  }

  if len(allowedHostNames) == 0 {
    return true
  }

  for index, element := range allowedHostNames {
    log.Debug("Checking index %d: %s against %s\n", index, element, input.Host)
    if element == input.Host {
      log.Debug("Match found on %d", index)
      return true
    }
  }

  return false
}
