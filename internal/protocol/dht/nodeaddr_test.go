package dht

import (
	"testing"

	"github.com/anacrolix/torrent/bencode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalNodeAddr(t *testing.T) {
	t.Parallel()

	var na NodeAddr

	require.NoError(t, na.UnmarshalBinary([]byte("\x01\x02\x03\x04\x05\x06")))
	assert.Equal(t, "1.2.3.4", na.IP.String())
}

func TestUnmarshalNodeAddrInvalidLength(t *testing.T) {
	t.Parallel()

	for _, b := range []string{"", "\x01", "\x01\x02\x03\x04\x05", "\x01\x02\x03\x04\x05\x06\x07"} {
		var na NodeAddr

		require.Error(t, na.UnmarshalBinary([]byte(b)))
	}
}

func TestUnmarshalMsgInvalidNodeAddr(t *testing.T) {
	t.Parallel()

	var m Msg

	require.Error(t, bencode.Unmarshal([]byte("d2:ip1:a1:t2:aa1:y1:re"), &m))
}
