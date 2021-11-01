# bunny-go
![CI](https://github.com/simplesurance/bunny-go/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/simplesurance/bunny-go)](https://goreportcard.com/report/github.com/simplesurance/bunny-go)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/simplesurance/bunny-go)

bunny-go is an unofficial Go package to interact with the [Bunny.net HTTP
API](https://docs.bunny.net/reference/bunnynet-api-overview). \
It aims to be a low-level API that represents Bunny API as close as possible.

## Features

The following [API
Endpoints](https://docs.bunny.net/reference/bunnynet-api-overview) are supported:

- [ ] bunny.net API
  - [ ] Billing
  - [ ] Stream Video Library
  - [ ] [Pull Zone](https://docs.bunny.net/reference/pullzonepublic_index)
    - [x] Add
    - [x] Update
    - [x] Delete
    - [x] Get
    - [x] List
    - [ ] Delete Edge Rule
    - [ ] Add/Update Edge Rule
    - [ ] Set Edge Rule Enabled
    - [ ] Get Statistics
    - [ ] Purge Cache
    - [ ] Load Free Certificate
    - [ ] Add Custom Certificate
    - [ ] Remove Certificate
    - [ ] Add Custom Hostname
    - [ ] Remove Custom Hostname
    - [ ] Set Force SSL
    - [ ] Reset Token Key
    - [ ] Add Allowed Referer
    - [ ] Remove Allowed Referer
    - [ ] Add Blocked Referer
    - [ ] Remove Blocked Referer
    - [ ] Add Blocked IP
    - [ ] Remove Blocked IP
  - [ ] Purge
  - [ ] Statistics
  - [ ] Storage Zone
  - [ ] User
- [ ] Edge Storage API
- [ ] Stream API

## Example

See [client_example_test.go](client_example_test.go)

## Design Principles

- URL parameters are always passed by value as method parameter.
- Data that is sent in the HTTP body is passed as struct
  pointer to API methods.
- Pointers instead of values are used to represent fields in body messages
  structs. This allows to represent unset fields correctly.
- Message field names should be as close as possible to the JSON message field
  names. Exception are permitted if the field in the JSON messages are
  inconsistent and different names are used in the API for the same setting.
  If names are inconsistent, the variant that is closer to the naming in the
  Bunny.Net Admin Panel should be chosen. The exception must be documented in
  the godoc.

## Status

The package is under initial development and should be considered as unstable.