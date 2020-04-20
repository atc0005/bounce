// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

module github.com/atc0005/bounce

go 1.13

// $ go list -m -versions github.com/apex/log
// github.com/apex/log v1.0.0 v1.1.0 v1.1.1 v1.1.2

// Use local copy of library package (instead of fetching remote content)
// replace github.com/atc0005/go-teams-notify => ../go-teams-notify
// replace github.com/atc0005/send2teams => ../send2teams

//
// require (
//
// ...
//  Note: Due to `replace` directive and `v0.0.0` here, we use the current
//  state of this library package from the local system instead of fetching
//  remote content
//	github.com/atc0005/go-teams-notify v0.0.0
//	github.com/atc0005/send2teams v0.0.0
// ...
//)

require (
	github.com/TylerBrock/colorjson v0.0.0-20180527164720-95ec53f28296
	github.com/apex/log v1.1.2

	//gopkg.in/dasrick/go-teams-notify.v1 v1.2.0

	// temporarily use our fork while developing changes for potential
	// inclusion in the upstream project
	github.com/atc0005/go-teams-notify v1.3.1-0.20200419155834-55cca556e726
	github.com/atc0005/send2teams v0.4.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/golang/gddo v0.0.0-20200324184333-3c2cc9a6329d
	github.com/hokaccha/go-prettyjson v0.0.0-20190818114111-108c894c2c0e // indirect
)
