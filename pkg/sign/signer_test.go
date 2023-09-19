package sign

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/valyala/bytebufferpool"
)

func TestSinger(t *testing.T) {
	bp := bytebufferpool.Get()
	defer bytebufferpool.Put(bp)
	sum := sha256.Sum256(bp.Bytes())
	payloadHash := hex.EncodeToString(sum[:])
	t.Log(payloadHash)
}
