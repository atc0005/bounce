// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Credit:

// https://golangcode.com/get-the-request-ip-addr/
// https://github.com/eddturtle/golangcode-site

package main

import (
	"net/http"

	"github.com/apex/log"
)

// GetIP gets a request's IP address by reading off the forwarded-for
// header (for proxies) and falls back to using the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	log.WithFields(log.Fields{
		"forwarded_header": forwarded,
	}).Debug("logging X-FORWARDED-FOR header")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
