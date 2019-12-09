package main

import (
	"fmt"
	"testing"
	"time"
)

func TestLicense(t *testing.T) {
	// 720 hours
	expire := time.Now().Add(30 * time.Hour).Unix()
	l := &License{
		Feature:         Features,
		ExpireTimestamp: expire,
	}
	enData := EncodeLicense(l)
	fmt.Println("encode data is : ", enData)
	license := DecodeLicense(enData)
	fmt.Printf("decode license is : %#v", license)
}
