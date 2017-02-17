package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_sessionTimeIDForTime(t *testing.T) {
	require.Equal(t, "20161220193456", sessionTimeIDForTime(time.Date(2016, 12, 20, 19, 34, 56, 0, time.UTC)))

	// non UTC time should still generate a UTC session time ID
	sameTimeInUTCPlusOne := "2016-12-20 20:34:56 +0100"
	nonUTCTime, err := time.Parse("2006-01-02 15:04:05 -0700", sameTimeInUTCPlusOne)
	require.NoError(t, err)
	require.Equal(t, "20161220193456", sessionTimeIDForTime(nonUTCTime))
}

func Test_readSessionStoreFromBytes(t *testing.T) {
	t.Log("Simple OK")
	{
		sessionStore, err := readSessionStoreFromBytes([]byte(`{"session_time_id": "20170215093215"}`))
		require.NoError(t, err)
		require.Equal(t, "20170215093215", sessionStore.SessionTimeID)
	}

	t.Log("Empty")
	{
		sessionStore, err := readSessionStoreFromBytes([]byte(`{}`))
		require.NoError(t, err)
		require.Equal(t, "", sessionStore.SessionTimeID)
	}
}
