// Copyright (c) 2018 Nikita Chisnikov
// Distributed under the MIT/X11 software license

package main

import (
	"log"
	"net"
	"bytes"
)

func getMacs() (macs []string) {
	ifcs, err := net.Interfaces()
	if err != nil {
		log.Println("getMacs():", err.Error())
		return macs
	}
	for i := range ifcs {
		if (ifcs[i].Flags&net.FlagUp==0)||(bytes.Compare(ifcs[i].HardwareAddr, nil)==0) {
			continue
		}
		macs = append(macs, ifcs[i].HardwareAddr.String())
	}
	return macs
}