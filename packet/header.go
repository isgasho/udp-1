package packet

type Header struct {
	SequenceNumber uint32
	Ack            Ack
}
