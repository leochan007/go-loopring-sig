package loopring

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"net/url"
	"strings"

	"github.com/loopring/go-loopring-sig/constants"
	"github.com/loopring/go-loopring-sig/eddsa"
)

func SignRequest(privateKey *eddsa.PrivateKey, method string, baseUrl string, path string, data string) (string, error) {
	methodParsed := strings.ToUpper(method)

	var params string
	if methodParsed == "GET" || methodParsed == "DELETE" {
		params = data
	} else if methodParsed == "POST" || methodParsed == "PUT" {
		params = url.QueryEscape(data)
	} else {
		return "", fmt.Errorf("%s is not supported yet", methodParsed)
	}

	uri := url.QueryEscape(fmt.Sprintf("%s%s", baseUrl, path))
	message := fmt.Sprintf(`%s&%s&%s`, methodParsed, uri, params)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(message)[:]))
	q, _ := new(big.Int).SetString(hash, 16)

	messageToSign := q.Mod(q, constants.Q)
	// messageToSignStr := messageToSign.String()

	sig := privateKey.SignPoseidon(messageToSign)

	sigReply := "0x" +
		fmt.Sprintf("%064s", sig.R8.X.Text(16)) +
		fmt.Sprintf("%064s", sig.R8.Y.Text(16)) +
		fmt.Sprintf("%064s", sig.S.Text(16))

	return sigReply, nil
}
