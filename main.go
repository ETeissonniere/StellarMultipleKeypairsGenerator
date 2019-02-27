package main

import (
    "log"
    "flag"

    "github.com/stellar/go/keypair"
)

var (
    publicKeyOnlyFile string
    allKeysFile string

    nbKeys uint
)

func newKeypair() (string, string) {
    kp, err := keypair.Random()
    if err != nil {
        log.Panic(err)
    }

    log.Println(kp.Address())

    return kp.Address(), kp.Seed()
}

func init() {
    flag.StringVar(&publicKeyOnlyFile, "pk-only", "public_keys.txt", "File with all the public keys")
    flag.StringVar(&allKeysFile, "keepass", "public_and_private_keys.keepass", "File with both the public and private keys")
    flag.UintVar(&nbKeys, "nb", 100, "Number of keys to generate")
}

func main() {
    flag.Parse()
}
