package packet

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Ack struct {
	sequenceNumber uint16
	acks           []uint16
}

func distance(newer, older uint16) uint16 {
	if older > newer {
		// Roll-over scenario.
		return math.MaxUint16 - older + newer
	}
	return newer - older
}

func (a *Ack) Serialize(payload []byte) ([]byte, error) {
	if len(payload) < 10 {
		return payload, fmt.Errorf("payload under 10 bytes. length: %d", len(payload))
	}
	binary.BigEndian.PutUint16(payload[0:], a.sequenceNumber)
	var acks uint64
	for _, ack := range a.acks {
		idx := distance(a.sequenceNumber, ack)
		acks |= (1 << idx)
	}
	binary.BigEndian.PutUint64(payload[2:], acks)
	return payload[6:], nil
}

func (a *Ack) Deserialize(payload []byte) ([]byte, error) {
	if len(payload) < 10 {
		return payload, fmt.Errorf("payload under 10 bytes. length: %d", len(payload))
	}
	a.sequenceNumber = binary.BigEndian.Uint16(payload[0:])
	acks := binary.BigEndian.Uint64(payload[2:])
	for i := 0; i < 64; i++ {
		if acks&(1<<i) != 0 {
			a.acks = append(a.acks, distance(a.sequenceNumber, uint16(i)))
		}
	}
	return payload, nil
}

func NewAck(SequenceNumber uint16, acks []uint16) (Ack, error) {
	for _, ack := range acks {
		if ack > SequenceNumber {
			return Ack{}, fmt.Errorf("ack > sequence number: %d vs %d", ack, SequenceNumber)
		}
		idx := distance(SequenceNumber, ack)
		if idx > 63 {
			return Ack{}, fmt.Errorf("distance too large too serialize: %d vs %d", SequenceNumber, ack)
		}
	}
	return Ack{
		sequenceNumber: SequenceNumber,
		acks:           acks,
	}, nil
}
