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

package tron

import (
	"github.com/blocktree/openwallet/log"
	"testing"
)

func TestGetAccountNet(t *testing.T) {

	var addr string
	addr = "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	accountNet, err := tw.GetAccountNet(addr)
	if err != nil {
		t.Errorf("GetAccountNet failed: %v\n", err)
		return
	}
	log.Infof("accountNet: %+v", accountNet)
}

func TestGetAccountResource(t *testing.T) {

	var addr string
	addr = "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	account, err := tw.GetAccountResource(addr)
	if err != nil {
		t.Errorf("GetAccountResource failed: %v\n", err)
		return
	}
	log.Infof("GetAccountResource: %+v", account)
}

func TestGetAccount(t *testing.T) {

	var addr string

	addr = "TF84rDKSaVXpFyzMvP8SPPzEChQ4QAP75N"

	if r, err := tw.GetAccount(addr); err != nil {
		t.Errorf("GetAccount failed: %v\n", err)
	} else {
		t.Logf("GetAccount return: \n\t%+v\n", r)
	}
}

func TestCreateAccount(t *testing.T) {

	var ownerAddress, accountAddress string = OWNERADDRESS, OWNERADDRESS

	if r, err := tw.CreateAccount(ownerAddress, accountAddress); err != nil {
		t.Errorf("CreateAccount failed: %v\n", err)
	} else {
		t.Logf("CreateAccount return: \n\t%+v\n", r)
	}
}

func TestUpdateAccount(t *testing.T) {

	var accountName, ownerAddress string = "XX2", OWNERADDRESS

	if r, err := tw.UpdateAccount(accountName, ownerAddress); err != nil {
		t.Errorf("UpdateAccount failed: %v\n", err)
	} else {
		t.Logf("UpdateAccount return: \n\t%+v\n", r)

		if r.Get("txID").String() == "" {
			t.Errorf("UpdateAccount failed: %v\n", "Data Invalid!")
		}
	}
}

func TestGetTRXAccount(t *testing.T) {

	var addr string

	addr = "TNmR1WCd2VD6vqXkmDjkDUPF8nXQAFmzJF"

	if r, _, err := tw.GetTRXAccount(addr); err != nil {
		t.Errorf("GetAccount failed: %v\n", err)
	} else {
		t.Logf("GetAccount return: \n\t%+v\n", r)
	}
}
