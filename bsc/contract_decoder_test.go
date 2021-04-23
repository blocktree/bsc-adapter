/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package bsc

import (
	"github.com/blocktree/openwallet/v2/openwallet"
	"testing"
)

func TestWalletManager_GetTokenBalanceByAddress(t *testing.T) {
	wm := testNewWalletManager()

	contract := openwallet.SmartContract{
		Address:  "0x4092678e4e78230f46a1534c0fbc8fa39780892b",
		Symbol:   "BSC",
		Name:     "OCoin",
		Token:    "OCN",
		Decimals: 18,
	}

	tokenBalances, err := wm.ContractDecoder.GetTokenBalanceByAddress(contract, "0xb45d2e41507cb95621d651d54645253a61b6b896", "0xfa23640bd91618d7fe79934d0eaaa747fb801fae")
	if err != nil {
		t.Errorf("GetTokenBalanceByAddress unexpected error: %v", err)
		return
	}
	for _, b := range tokenBalances {
		t.Logf("token balance: %+v", b.Balance)
	}
}
