package main

type Profile struct {
	Name   string `datastore:"name, noindex"`
	Email  string `datastore:"email"`
	Gender int `datastore:"gender, noindex"`
}