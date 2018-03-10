package chanmanager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClaim(t *testing.T) {
	chm := New()
	set := map[[2]byte]bool{}
	for i := 0; i < 256*256; i++ {
		id, channel, err := chm.Claim()
		assert.NoError(t, err)
		assert.NotNil(t, channel)
		assert.False(t, set[id], "Id %s was already taken", id)
		set[id] = true
	}
	_, _, err := chm.Claim()
	assert.Error(t, err)
}

func TestClaim_AfterUnclaim(t *testing.T) {
	chm := New()
	set := map[[2]byte]bool{}
	for i := 0; i < 256*256; i++ {
		id, channel, err := chm.Claim()
		assert.NoError(t, err)
		assert.NotNil(t, channel)
		assert.False(t, set[id], "Id %s was already taken", id)
		set[id] = true
	}
	expectedID := [2]byte{13, 14}
	chm.Unclaim(expectedID)

	actualID, channel, err := chm.Claim()
	assert.NoError(t, err)
	assert.NotNil(t, channel)
	assert.Equal(t, expectedID, actualID)
}

func TestClaimWithID_WhenIDIsFree(t *testing.T) {
	chm := New()

	channel, err := chm.ClaimWithID([2]byte{13, 14})

	assert.NoError(t, err)
	assert.NotNil(t, channel)
}

func TestClaimWithID_WhenIDIsTakenByClaimWithID(t *testing.T) {
	chm := New()
	chm.ClaimWithID([2]byte{13, 14})

	_, err := chm.ClaimWithID([2]byte{13, 14})

	assert.Error(t, err)
}

func TestClaimWithID_WhenIDIsTakenByClaim(t *testing.T) {
	chm := New()
	id, _, _ := chm.Claim()

	_, err := chm.ClaimWithID(id)

	assert.Error(t, err)
}

func TestClaim_WhenIDIsTakenByClaimWithID(t *testing.T) {
	chm := New()
	takenID := [2]byte{0, 0}
	chm.ClaimWithID(takenID)

	id, channel, err := chm.Claim()

	assert.NoError(t, err)
	assert.NotNil(t, channel)
	assert.NotEqual(t, takenID, id)
}

func TestSendData_WhenIDIsTaken(t *testing.T) {
	chm := New()
	takenID := [2]byte{0, 0}
	passedValue := struct{}{}

	channel, _ := chm.ClaimWithID(takenID)
	err := chm.SendData(takenID, passedValue)

	assert.NoError(t, err)
	assert.Equal(t, passedValue, <-channel)
}

func TestSendData_WhenIDIsNotTaken(t *testing.T) {
	chm := New()
	notTakenID := [2]byte{0, 0}

	err := chm.SendData(notTakenID, struct{}{})

	assert.Error(t, err)
}

func BenchmarkClaim(b *testing.B) {
	chm := New()
	for i := 0; i < b.N; i++ {
		id, _, err := chm.Claim()
		if err != nil {
			b.Error(err)
		}
		chm.Unclaim(id)
	}
}
