package dispatcher

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	ErrClosed = fmt.Errorf("connection closed")
)

type SessionWriter struct {
	device *Device

	id        string
	dataChan  chan []byte
	closeChan chan struct{}
	closeOnce sync.Once
}

func NewSessionWriter(device *Device) *SessionWriter {
	return &SessionWriter{
		device:    device,
		id:        uuid.New().String(),
		dataChan:  make(chan []byte, 5),
		closeChan: make(chan struct{}),
	}
}

func (sw *SessionWriter) Write(data []byte) error {
	return sw.device.writeMsgToDevice(sw.id, data)
}

func (sw *SessionWriter) Read() ([]byte, error) {
	select {
	case data := <-sw.dataChan:
		return data, nil
	case <-sw.closeChan:
		return nil, ErrClosed
	}
}

func (sw *SessionWriter) Close() {
	sw.closeOnce.Do(func() {
		sw.device.removeSessionWriter(sw.id)
		close(sw.closeChan)
	})
}

func (sw *SessionWriter) fwToSession(data []byte) {
	select {
	case sw.dataChan <- data:
	default:
		log.Errorf("messaged dropped")
	}
	return
}
