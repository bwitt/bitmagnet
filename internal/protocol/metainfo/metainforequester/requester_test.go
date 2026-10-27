package metainforequester

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pieceMessage(piece int, data []byte) []byte {
	payload := append([]byte(fmt.Sprintf("\x14\x01d8:msg_typei1e5:piecei%dee", piece)), data...)
	msg := binary.BigEndian.AppendUint32(nil, uint32(len(payload)))

	return append(msg, payload...)
}

func TestReadAllPieces(t *testing.T) {
	t.Parallel()

	full := bytes.Repeat([]byte{1}, 16*1024)
	last := bytes.Repeat([]byte{2}, 100)
	metadataSize := uint(len(full) + len(last))

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		r := bytes.NewReader(append(pieceMessage(0, full), pieceMessage(1, last)...))
		result, err := readAllPieces(r, metadataSize)
		require.NoError(t, err)
		assert.Equal(t, append(full, last...), result)
	})

	for name, msgs := range map[string][]byte{
		"negative index":     pieceMessage(-1, full),
		"out of range index": pieceMessage(2, full),
		"duplicate piece":    append(pieceMessage(0, full), pieceMessage(0, full)...),
		"overflowing piece":  pieceMessage(1, full),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := readAllPieces(bytes.NewReader(msgs), metadataSize)
			require.Error(t, err)
		})
	}
}
