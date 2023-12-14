package dispatcher

type NetConn interface {
	Write(id string, buf []byte) error
	Read() ([]byte, string, error)
}
