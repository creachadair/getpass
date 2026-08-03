# getpass

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/creachadair/getpass)
[![CI](https://github.com/creachadair/getpass/actions/workflows/go-presubmit.yml/badge.svg?event=push&branch=main)](https://github.com/creachadair/getpass/actions/workflows/go-presubmit.yml)

This repository provides a Go package to read passwords from the terminal with
echo disabled. This implementation relies on the [x/sys/unix][unix] package to
read and write terminal settings.

[unix]: http://godoc.org/golang.org/x/sys/unix
