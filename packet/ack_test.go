package packet

import (
	"log"
	"testing"
)

func ackEquals(a, b Ack) bool {
	if a.sequenceNumber != b.sequenceNumber {
		return false
	}
	if len(a.acks) != len(b.acks) {
		return false
	}
	acks := make(map[uint16]struct{})
	for _, ack := range a.acks {
		acks[ack] = struct{}{}
	}
	for _, ack := range b.acks {
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
			ack:  Ack{42, []uint16{41, 40, 36}},
		},
		{
			name: "rollover",
			ack:  Ack{8, []uint16{65534, 65533}},
		},
		{
			name: "rollover pt. 2",
			ack:  Ack{50, []uint16{49, 38, 36, 65534, 65533}},
		},
		{
			name: "distance just right",
			ack:  Ack{100, []uint16{37}},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			payload := make([]byte, 1200)
			_, err := tc.ack.Serialize(payload)
			if err != nil {
				t.Fatal("unexpected erorr", err)
			}

			var deserialized Ack
			_, err = deserialized.Deserialize(payload)
			if err != nil {
				t.Fatal("unexpected erorr", err)
			}

			log.Println("serialized  ", tc.ack)
			log.Println("deserialized", deserialized)
			if !ackEquals(tc.ack, deserialized) {
				t.Error("acks not equal after serialize -> deserialize")
			}
		})
	}
}
