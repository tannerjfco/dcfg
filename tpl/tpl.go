package tpl

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/drud/drud-go/utils/stringutil"
)

// Tpl is the interface that each plugin must implement
type Tpl interface {
	WriteConfig(in *Config) error
	PlaceFiles(in *Config, move bool) error
	WebConfig(in *Config) error
}

// TplMap is used to retrieve the correct plugin
var TplMap = map[string]Tpl{
	"drupal":    &DrupalConfig{},
	"wordpress": &WordpressConfig{},
}

// PassTheSalt generates a hash salt
func PassTheSalt() string {
	salt := sha256.New()
	random := stringutil.RandomString(20)
	salt.Write([]byte(random))

	return hex.EncodeToString(salt.Sum(nil))
}

// SlashIt ensures you have a preceding or trailing slash on a string
func SlashIt(val string, trailing bool) string {
	if trailing && !strings.HasSuffix(val, "/") {
		return val + "/"
	} else if !trailing && !strings.HasPrefix(val, "/") {
		return "/" + val
	}
	return val
}
