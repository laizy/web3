package codec

import "encoding/binary"

type Serializable interface {
	Serialization(sink *ZeroCopySink)
}

func SerializeToBytes(values ...Serializable) []byte {
	sink := NewZeroCopySink(nil)
	for _, val := range values {
		val.Serialization(sink)
	}

	return sink.Bytes()
}

func BytesLEToUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func Uint16ToBytesLE(val uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, val)
	return data
}

func BytesLEToUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func Uint32ToBytesLE(val uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, val)
	return data
}

func BytesLEToUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func Uint64ToBytesLE(val uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, val)
	return data
}
