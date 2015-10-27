# Couch

[![Build Status](https://travis-ci.org/fcavani/couch.svg?branch=master)](https://travis-ci.org/fcavani/couch) [![GoDoc](https://godoc.org/github.com/fcavani/couch?status.svg)](https://godoc.org/github.com/fcavani/couch)

Couch package implements another driver for Couch databases.
Here you can create, view and delete a database. The three
basic function are suported: put, get, delete and nothing 
more. Structs will be serialized with the json package of 
Go and need to have ``json:"_id"`` for the id field and
``json:"_rev,omitempty"`` for the revision field.

# Examples

See the couch_test.go.

# Install
go get github.com/fcavani/couch
