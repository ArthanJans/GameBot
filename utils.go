package main

import "strings"

func idFromTag(tag string) string {
	return strings.TrimPrefix(strings.TrimSuffix(tag, ">"), "<@")
}
