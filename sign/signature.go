package sign

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"strconv"
)

func CreateSignature(text string, timestamp int64, prevSignature []byte) []byte {
	passphrase, _, privateKey := getCredentials()
	message := crypto.NewPlainMessageFromString(text + strconv.FormatInt(timestamp, 10) + string(prevSignature))

	privateKeyObj, _ := crypto.NewKeyFromArmored(privateKey)
	unlockedKeyObj, _ := privateKeyObj.Unlock([]byte(passphrase))
	signingKeyRing, _ := crypto.NewKeyRing(unlockedKeyObj)
	pgpSignature, _ := signingKeyRing.SignDetached(message)

	return pgpSignature.Data
}

func VerifySignature(text string, timestamp int64, prevSignature []byte, signature []byte) bool {
	_, publicKey, _ := getCredentials()
	message := crypto.NewPlainMessageFromString(text + strconv.FormatInt(timestamp, 10) + string(prevSignature))

	pgpSignature := &crypto.PGPSignature{Data: signature}
	publicKeyObj, _ := crypto.NewKeyFromArmored(publicKey)
	signingKeyRing, _ := crypto.NewKeyRing(publicKeyObj)
	err := signingKeyRing.VerifyDetached(message, pgpSignature, crypto.GetUnixTime())

	return err == nil
}
