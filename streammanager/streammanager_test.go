package streammanager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClaim(t *testing.T) {
	sm := New()
	set := map[StreamID]bool{}
	for i := 0; i < 256*256; i++ {
		id, channel, err := sm.Claim()
		assert.NoError(t, err)
		assert.NotNil(t, channel)
		assert.False(t, set[id], "Id %s was already taken", id)
		set[id] = true
	}
	_, _, err := sm.Claim()
	assert.Error(t, err)
}

func TestClaim_AfterUnclaim(t *testing.T) {
	sm := New()
	set := map[StreamID]bool{}
	for i := 0; i < 256*256; i++ {
		id, channel, err := sm.Claim()
		assert.NoError(t, err)
		assert.NotNil(t, channel)
		assert.False(t, set[id], "Id %s was already taken", id)
		set[id] = true
	}
	expectedID := StreamID{13, 14}
	sm.Unclaim(expectedID)

	actualID, channel, err := sm.Claim()
	assert.NoError(t, err)
	assert.NotNil(t, channel)
	assert.Equal(t, expectedID, actualID)
}

func TestClaimWithID_WhenIDIsFree(t *testing.T) {
	sm := New()

	channel, err := sm.ClaimWithID(StreamID{13, 14})

	assert.NoError(t, err)
	assert.NotNil(t, channel)
}

func TestClaimWithID_WhenIDIsTakenByClaimWithID(t *testing.T) {
	sm := New()
	sm.ClaimWithID(StreamID{13, 14})

	_, err := sm.ClaimWithID(StreamID{13, 14})

	assert.Error(t, err)
}

func TestClaimWithID_WhenIDIsTakenByClaim(t *testing.T) {
	sm := New()
	id, _, _ := sm.Claim()

	_, err := sm.ClaimWithID(id)

	assert.Error(t, err)
}

func TestClaim_WhenIDIsTakenByClaimWithID(t *testing.T) {
	sm := New()
	takenID := StreamID{0, 0}
	sm.ClaimWithID(takenID)

	id, channel, err := sm.Claim()

	assert.NoError(t, err)
	assert.NotNil(t, channel)
	assert.NotEqual(t, takenID, id)
}

func TestSendData_WhenIDIsTaken(t *testing.T) {
	sm := New()
	takenID := StreamID{0, 0}
	passedValue := &ServerResponse{}

	channel, _ := sm.ClaimWithID(takenID)
	err := sm.SendData(takenID, passedValue)

	assert.NoError(t, err)
	assert.Equal(t, passedValue, <-channel)
}

func TestSendData_WhenIDIsNotTaken(t *testing.T) {
	sm := New()
	notTakenID := StreamID{0, 0}

	err := sm.SendData(notTakenID, &ServerResponse{})

	assert.Error(t, err)
}

func BenchmarkClaim(b *testing.B) {
	sm := New()
	for i := 0; i < b.N; i++ {
		id, _, err := sm.Claim()
		if err != nil {
			b.Error(err)
		}
		sm.Unclaim(id)
	}
}
