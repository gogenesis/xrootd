package chanmanager

import (
	"github.com/pkg/errors"
	"sync"
)

// Chanmanager manages channels by their ids. Basically, it's a map[[2]byte] chan<-interface{} with methods
// to claim, free and pass data to some channel by id.
type Chanmanager struct {
	dataWaiters map[[2]byte]chan<- interface{}
	mutex       sync.Mutex
	freeIds     chan [2]byte
}

// New creates new Chanmanager
func New() *Chanmanager {
	chm := Chanmanager{make(map[[2]byte]chan<- interface{}), sync.Mutex{}, make(chan [2]byte, 256*256)}

	for firstByte := uint16(0); firstByte < 256; firstByte++ {
		for secondByte := uint16(0); secondByte < 256; secondByte++ {
			chm.freeIds <- [2]byte{byte(firstByte), byte(secondByte)}
		}
	}

	return &chm
}

// Claim searches for unclaimed id and returns corresponding channel
func (chm *Chanmanager) Claim() (id [2]byte, channel <-chan interface{}, err error) {
	bidirectionalChannel := make(chan interface{}, 1)
	chm.mutex.Lock()
	defer chm.mutex.Unlock()

	for channel == nil && err == nil {
		select {
		case id = <-chm.freeIds:
			if _, ok := chm.dataWaiters[id]; ok { // Skip id if it was already taken manually via ClaimWithID
				continue
			}

			chm.dataWaiters[id] = bidirectionalChannel
			channel = bidirectionalChannel
		default:
			err = errors.New("Cannot obtain free channel")
		}
	}
	return
}

// ClaimWithID checks if id is unclaimed and returns corresponding channel in case of success
func (chm *Chanmanager) ClaimWithID(id [2]byte) (channel <-chan interface{}, err error) {
	bidirectionalChannel := make(chan interface{}, 1)
	chm.mutex.Lock()
	defer chm.mutex.Unlock()
	if _, ok := chm.dataWaiters[id]; ok {
		err = errors.Errorf("Channel with id %s is already taken", id)
	} else {
		chm.dataWaiters[id] = bidirectionalChannel
		channel = bidirectionalChannel
	}
	return
}

// Unclaim marks channel with specified id as unclaimed
func (chm *Chanmanager) Unclaim(id [2]byte) {
	chm.mutex.Lock()
	delete(chm.dataWaiters, id)
	chm.mutex.Unlock()
	chm.freeIds <- id
}

// SendData sends data to channel with specific id
func (chm *Chanmanager) SendData(id [2]byte, data interface{}) error {
	chm.mutex.Lock()
	defer chm.mutex.Unlock()

	if _, ok := chm.dataWaiters[id]; !ok {
		return errors.Errorf("Cannot find data waiter for id %s", id)
	}

	chm.dataWaiters[id] <- data

	return nil
}
