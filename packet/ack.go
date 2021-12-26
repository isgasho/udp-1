package packet

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Ack struct {
	SequenceNumber uint16
	Acks           []uint16
}

func distance(newer, older uint16) uint16 {
	if older > newer {
		// Roll-over scenario.
		return math.MaxUint16 - older + newer
	}
	return newer - older
}

func (a *Ack) Serialize(payload []byte) ([]byte, error) {
	binary.BigEndian.PutUint16(payload[0:], a.SequenceNumber)
	var acks uint64
	for _, ack := range a.Acks {
		idx := distance(a.SequenceNumber, ack)
		if idx > 63 {
			return nil, fmt.Errorf("distance too large too serialize: %d vs %d", a.SequenceNumber, ack)
		}

		acks |= (1 << idx)
	}
	binary.BigEndian.PutUint64(payload[2:], acks)
	return payload[6:], nil
}

func (a *Ack) Deserialize(payload []byte) error {
	if len(payload) < 8 {
		return fmt.Errorf("payload under 8 bytes. length: %d", len(payload))
	}

	a.SequenceNumber = binary.BigEndian.Uint16(payload[0:])
	acks := binary.BigEndian.Uint64(payload[2:])
	for i := 0; i < 64; i++ {
		if acks&(1<<i) != 0 {
			a.Acks = append(a.Acks, distance(a.SequenceNumber, uint16(i)))
		}
	}
	return nil
}

func NewAck(SequenceNumber uint16, acks []uint16) Ack {
	return Ack{
		SequenceNumber: SequenceNumber,
		Acks:           acks,
	}
}
