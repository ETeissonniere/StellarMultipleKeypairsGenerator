package main

import (
    "log"

    "github.com/stellar/go/keypair"
)

func newKeypair() (string, string) {
    kp, err := keypair.Random()
    if err != nil {
        check(err)
    }

    log.Println(kp.Address())

    return kp.Address(), kp.Seed()
}
