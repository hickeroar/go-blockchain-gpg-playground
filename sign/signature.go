package sign

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// We would want error checking every step of the way instead of the ignores (underscores).

func CreateSignature(text string, hash []byte) []byte {
	passphrase, _, privateKey := getCredentials()
	message := crypto.NewPlainMessageFromString(text + string(hash))

	privateKeyObj, _ := crypto.NewKeyFromArmored(privateKey)
	unlockedKeyObj, _ := privateKeyObj.Unlock([]byte(passphrase))
	signingKeyRing, _ := crypto.NewKeyRing(unlockedKeyObj)
	pgpSignature, _ := signingKeyRing.SignDetached(message)

	return pgpSignature.Data
}

func VerifySignature(text string, hash []byte, signature []byte) bool {
	_, publicKey, _ := getCredentials()
	message := crypto.NewPlainMessageFromString(text + string(hash))

	pgpSignature := &crypto.PGPSignature{Data: signature}
	publicKeyObj, _ := crypto.NewKeyFromArmored(publicKey)
	signingKeyRing, _ := crypto.NewKeyRing(publicKeyObj)
	err := signingKeyRing.VerifyDetached(message, pgpSignature, crypto.GetUnixTime())

	return err == nil
}
