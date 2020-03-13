package impl

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(sha256Str([]byte("\"\"+abc")))
}

func sha256ForTest(x []byte) string {
	// y := sha256.Sum256(x)
	// res := hex.EncodeToString(y[:])
	// bs, _ := hex.DecodeString(res)
	return ""
}
