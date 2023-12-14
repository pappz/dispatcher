package dispatcher

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type Device struct {
	id   string
	conn NetConn

	sessionWriters map[string]*SessionWriter
	swLock         sync.RWMutex
}

func NewDevice(id string, conn NetConn) *Device {
	d := &Device{
		id:             id,
		conn:           conn,
		sessionWriters: make(map[string]*SessionWriter),
	}

	go d.startReader()
	return d
}

func (d *Device) OpenNewSession() *SessionWriter {
	d.swLock.Lock()
	defer d.swLock.Unlock()

	sw := NewSessionWriter(d)
	d.sessionWriters[sw.id] = sw
	return sw
}

func (d *Device) startReader() {
	for {
		data, sessionId, err := d.conn.Read()
		if err != nil {
			log.Infof("stop dispatcher reader for device: %s", d.id)
			return
		}

		d.swLock.RLock()
		sw, ok := d.sessionWriters[sessionId]
		d.swLock.RUnlock()
		if !ok {
			log.Errorf("session writer not found, message dropped")
			continue
		}
		sw.fwToSession(data)
	}
}

func (d *Device) writeMsgToDevice(id string, data []byte) error {
	return d.conn.Write(id, data)
}

func (d *Device) removeSessionWriter(id string) {
	d.swLock.Lock()
	defer d.swLock.Unlock()

	delete(d.sessionWriters, id)

}
