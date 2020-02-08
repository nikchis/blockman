// Copyright (c) 2018 Nikita Chisnikov
// Distributed under the MIT/X11 software license

package main

type StakingStatus struct {
	ValidTime       bool `json:"validtime"`
	HaveConnections bool `json:"haveconnections"`
	Unlocked        bool `json:"walletunlocked"`
	MinTabCoins     bool `json:"mintablecoins"`
	EnoughCoins     bool `json:"enoughcoins"`
	MnSync          bool `json:"mnsync"`
	Staking         bool `json:"staking status"`
}

type UnlockError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
