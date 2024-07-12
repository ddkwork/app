package main

import (
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	t.Skip()
	g := stream.NewGeneratedFile()
	g.Enum("Crypt", []string{
		"Symmetry",
		"Asymmetrical",
		"Hash",
		"Encoding",
		"Tool",
	}, nil)
	g.Enum("CryptName", []string{
		"Aes",
		"Des",
		"Des3 ",
		"Tea",
		"Blowfish",
		"TwoFish",
		"Rc4",
		"Rc2",

		"Rsa",
		"Ecc",
		"Dsa",
		"Pgp",
		"Sm4",
		"Sm2",

		"Hmac",
		"hashAll",

		"Base64",
		"Base32",
		"Gzip",

		"trimSpace",
		"swap",
		"request header",
		"timeStamp",
	}, nil)
}
