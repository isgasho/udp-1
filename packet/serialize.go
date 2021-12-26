package packet

type Serialize interface {
	Serialize([]byte) ([]byte, error)
}

type Deserialize interface {
	Deserialize([]byte) ([]byte, error)
}
