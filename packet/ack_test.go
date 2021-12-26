package packet

import (
	"log"
	"testing"
)

func ackEquals(a, b Ack) bool {
	if a.SequenceNumber != b.SequenceNumber {
		return false
	}
	if len(a.Acks) != len(b.Acks) {
		return false
	}
	acks := make(map[uint16]struct{})
	for _, ack := range a.Acks {
		acks[ack] = struct{}{}
	}
	for _, ack := range b.Acks {
		if _, ok := acks[ack]; !ok {
			return false
		}
	}
	return true
}

func TestAckSerializeDeserialize(t *testing.T) {
	tcs := []struct {
		name      string
		ack       Ack
		expectErr bool
	}{
		{
			name: "basic",
			ack:  NewAck(42, []uint16{41, 40, 36}),
		},
		{
			name: "rollover",
			ack:  NewAck(8, []uint16{65534, 65533}),
		},
		{
			name: "rollover pt. 2",
			ack:  NewAck(50, []uint16{65534, 65533}),
		},
		{
			name: "distance just right",
			ack:  NewAck(100, []uint16{37}),
		},
		{
			name:      "distance too large",
			ack:       NewAck(100, []uint16{36}),
			expectErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			payload := make([]byte, 1200)
			_, err := tc.ack.Serialize(payload)
			if tc.expectErr {
				log.Println(err)
				if err == nil {
					t.Error("expected error")
				}
				return
			} else {
				if err != nil {
					t.Error("unexpected erorr", err)
					return
				}
			}

			var deserialized Ack
			deserialized.Deserialize(payload)

			log.Println("serialized", tc.ack)
			log.Println("deserialized", deserialized)
			if !ackEquals(tc.ack, deserialized) {
				t.Error("acks not equal after serialize -> deserialize")
			}
		})
	}
}
