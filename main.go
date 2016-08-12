// Package modulist is a web-based review tool for module
// descriptions developed by Freitagsrunde at TU Berlin.
package main

import (
	"fmt"
)

func main() {

	// Init app.
	app := InitApp()

	// Run MODULIST either with or without TLS.
	if app.TLS {
		app.Router.RunTLS(fmt.Sprintf("%s:%s", app.IP, app.Port), app.TLSCertFile, app.TLSKeyFile)
	} else {
		app.Router.Run(fmt.Sprintf("%s:%s", app.IP, app.Port))
	}
}
