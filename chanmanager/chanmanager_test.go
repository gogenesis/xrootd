package chanmanager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClaimID(t *testing.T) {
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

func TestClaimIDAfterUnclaim(t *testing.T) {
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

func TestClaimWithIDWhenIDIsFree(t *testing.T) {
	chm := New()

	channel, err := chm.ClaimWithID([2]byte{13, 14})

	assert.NoError(t, err)
	assert.NotNil(t, channel)
}

func TestClaimWithIDWhenIDIsTakenByClaimWithID(t *testing.T) {
	chm := New()
	chm.ClaimWithID([2]byte{13, 14})

	_, err := chm.ClaimWithID([2]byte{13, 14})

	assert.Error(t, err)
}

func TestClaimWithIDWhenIDIsTakenByClaim(t *testing.T) {
	chm := New()
	id, _, _ := chm.Claim()

	_, err := chm.ClaimWithID(id)

	assert.Error(t, err)
}

func TestClaimWhenIDIsTakenByClaimWithID(t *testing.T) {
	chm := New()
	takenID := [2]byte{0, 0}
	chm.ClaimWithID(takenID)

	id, channel, err := chm.Claim()

	assert.NoError(t, err)
	assert.NotNil(t, channel)
	assert.NotEqual(t, takenID, id)
}

func TestSendDataWhenIDIsTaken(t *testing.T) {
	chm := New()
	takenID := [2]byte{0, 0}
	passedValue := struct{}{}

	channel, _ := chm.ClaimWithID(takenID)
	err := chm.SendData(takenID, passedValue)

	assert.NoError(t, err)
	assert.Equal(t, passedValue, <-channel)
}

func TestSendDataWhenIDIsNotTaken(t *testing.T) {
	chm := New()
	notTakenID := [2]byte{0, 0}

	err := chm.SendData(notTakenID, struct{}{})

	assert.Error(t, err)
}

func BenchmarkClaimID(b *testing.B) {
	chm := New()
	for i := 0; i < 256*256; i++ {
		_, _, err := chm.Claim()
		assert.NoError(b, err)
	}
}
