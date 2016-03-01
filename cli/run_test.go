package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LogBuffer_ReadRunes(t *testing.T) {
	t.Log("Read max 0 runes")
	{
		logbuff := LogBuffer{}
		str, isEOF := logbuff.ReadRunes(0)
		require.Equal(t, "", str)
		require.Equal(t, false, isEOF)
	}

	t.Log("Read from empty buffer")
	{
		logbuff := LogBuffer{}
		str, isEOF := logbuff.ReadRunes(100)
		require.Equal(t, "", str)
		require.Equal(t, true, isEOF)
	}

	t.Log("Read from a simple buffer - max runes count > buffer size")
	{
		logbuff := LogBuffer{}
		_, err := logbuff.Write([]byte("0123456789"))
		require.NoError(t, err)
		//
		str, isEOF := logbuff.ReadRunes(1000)
		require.Equal(t, "0123456789", str)
		require.Equal(t, true, isEOF)
	}

	t.Log("Read chunks from buffer")
	{
		logbuff := LogBuffer{}
		_, err := logbuff.Write([]byte("111111111122222222223333333333"))
		require.NoError(t, err)
		// read first chunk
		str, isEOF := logbuff.ReadRunes(10)
		require.Equal(t, "1111111111", str)
		require.Equal(t, false, isEOF)
		// read second chunk
		str, isEOF = logbuff.ReadRunes(10)
		require.Equal(t, "2222222222", str)
		require.Equal(t, false, isEOF)
		// read third/last chunk
		str, isEOF = logbuff.ReadRunes(10)
		require.Equal(t, "3333333333", str)
		require.Equal(t, false, isEOF)
		// EOF - no more data in the buffer
		str, isEOF = logbuff.ReadRunes(10)
		require.Equal(t, "", str)
		require.Equal(t, true, isEOF)
	}
}
