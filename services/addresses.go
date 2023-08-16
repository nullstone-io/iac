package services

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type AddressLoader string

func (a AddressLoader) Get() string {
	if val := os.Getenv(strings.ToUpper(fmt.Sprintf("%s_addr", a))); val != "" {
		return val
	}
	return (&url.URL{Scheme: "http", Host: QualifySubdomain(string(a))}).String()
}

var (
	EnigmaAddress = AddressLoader("enigma")
	FurionAddress = AddressLoader("furion")
	VoidAddress   = AddressLoader("void")
)
