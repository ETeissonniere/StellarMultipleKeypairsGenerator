package main

import (
    "os"
    "log"
    "flag"
    "time"
    "strings"

    "math/rand"

    "github.com/tobischo/gokeepasslib/v2"
)

var (
    publicKeyOnlyFile string
    allKeysFile string

    nbKeys int
)

func generatePassword() string {
    chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789#!;,/@()[]{}$?ù%çà;.:/?=+")
    length := 20

    var b strings.Builder

    for i := 0; i < length; i++ {
        b.WriteRune(chars[rand.Intn(len(chars))])
    }

    return b.String()
}

func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}}
}

func mkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value, Protected: true}}
}

func init() {
    rand.Seed(time.Now().UnixNano())

    flag.StringVar(&publicKeyOnlyFile, "pk-only", "public_keys.txt", "File with all the public keys")
    flag.StringVar(&allKeysFile, "keepass", "public_and_private_keys.kdbx", "File with both the public and private keys")
    flag.IntVar(&nbKeys, "nb", 100, "Number of keys to generate")
}

func main() {
    flag.Parse()

    kdbFile, err := os.Create(allKeysFile)
    check(err)
    defer kdbFile.Close()

    kdbPass := generatePassword()

    kdbGroup := gokeepasslib.NewGroup()
    kdbGroup.Name = "Stellar keys"

    for i := 0; i < nbKeys; i++ {
        pk, sk := newKeypair()

        entry := gokeepasslib.NewEntry()
        entry.Values = append(entry.Values, mkValue("Title", pk))
        entry.Values = append(entry.Values, mkProtectedValue("Private key", sk))

        kdbGroup.Entries = append(kdbGroup.Entries, entry)
    }

    kdb := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(kdbPass),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{kdbGroup},
			},
		},
	}

	// Lock entries using stream cipher
	kdb.LockProtectedEntries()

	// and encode it into the file
	kdbEncoder := gokeepasslib.NewEncoder(kdbFile)
    err = kdbEncoder.Encode(kdb)
    check(err)

    log.Printf("Keepass password: %s\n", kdbPass)
}
