package internal

import (
	"encoding/binary"
	"fmt"
)

// Protocol messages
const (
	MessageTypeProduced int16 = iota
	MessageTypeFetch
	MessageTypeMetadata
	MessageTypeCommitOffset
	MessageTypeFetchOffset
)

// SerializeMessage serializes a message for transmission
func SerializeMessage(msgType int16, data []byte) ([]byte, error) {
	// Format: [type:2][length:4][data]
	result := make([]byte, 6+len(data))

	binary.BigEndian.PutUint16(result[0:2], uint16(msgType))
	binary.BigEndian.PutUint32(result[2:6], uint32(len(data)))
	copy(result[6:], data)

	return result, nil
}

// DeserializeMessage deserializes a received message
func DeserializeMessage(data []byte) (int16, []byte, error) {
	if len(data) < 6 {
		return 0, nil, fmt.Errorf("message too short: %d bytes", len(data))
	}

	msgType := int16(binary.BigEndian.Uint16(data[0:2]))
	length := binary.BigEndian.Uint32(data[2:6])

	if uint32(len(data)-6) < length {
		return 0, nil, fmt.Errorf("incomplete message: expected %d bytes, got %d",
			length, len(data)-6)
	}

	return msgType, data[6 : 6+length], nil
}

// SerializeProduceRequest serializes a produce request
func SerializeProduceRequest(topic string, partition int32, key, value []byte) []byte {
	// Format: [topic_len:2][topic][partition:4][key_len:4][key][value_len:4][value]
	buf := make([]byte, 0, 2+len(topic)+4+4+len(key)+4+len(value))

	// Topic
	buf = append(buf, byte(len(topic)>>8), byte(len(topic)))
	buf = append(buf, []byte(topic)...)

	// Partition
	partBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(partBuf, uint32(partition))
	buf = append(buf, partBuf...)

	// Key
	keyLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(keyLenBuf, uint32(len(key)))
	buf = append(buf, keyLenBuf...)
	buf = append(buf, key...)

	// Value
	valLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(valLenBuf, uint32(len(value)))
	buf = append(buf, valLenBuf...)
	buf = append(buf, value...)

	return buf
}

// SerializeFetchRequest serializes a fetch request
func SerializeFetchRequest(topic string, partition int32, offset int64) []byte {
	// Format: [topic_len:2][topic][partition:4][offset:8]
	buf := make([]byte, 0, 2+len(topic)+4+8)

	// Topic
	buf = append(buf, byte(len(topic)>>8), byte(len(topic)))
	buf = append(buf, []byte(topic)...)

	// Partition
	partBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(partBuf, uint32(partition))
	buf = append(buf, partBuf...)

	// Offset
	offsetBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(offsetBuf, uint64(offset))
	buf = append(buf, offsetBuf...)

	return buf
}

// DeserializeProduceResponse unpacks a produce response
func DeserializeProduceResponse(data []byte) (partition int32, offset int64, errCode int16, err error) {
	if len(data) < 14 {
		return 0, 0, 0, fmt.Errorf("response too short: %d bytes", len(data))
	}

	partition = int32(binary.BigEndian.Uint32(data[0:4]))
	offset = int64(binary.BigEndian.Uint64(data[4:12]))
	errCode = int16(binary.BigEndian.Uint16(data[12:14]))

	return partition, offset, errCode, nil
}

// DeserializeFetchResponse unpacks a fetch response
func DeserializeFetchResponse(data []byte) (key, value []byte, errCode int16, err error) {
	if len(data) < 10 {
		return nil, nil, 0, fmt.Errorf("response too short: %d bytes", len(data))
	}

	keyLen := binary.BigEndian.Uint32(data[0:4])
	offset := 4

	if offset+int(keyLen) > len(data) {
		return nil, nil, 0, fmt.Errorf("incomplete key in response")
	}

	key = data[offset : offset+int(keyLen)]
	offset += int(keyLen)

	if offset+4 > len(data) {
		return nil, nil, 0, fmt.Errorf("incomplete value length in response")
	}

	valLen := binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	if offset+int(valLen) > len(data) {
		return nil, nil, 0, fmt.Errorf("incomplete value in response")
	}

	value = data[offset : offset+int(valLen)]
	offset += int(valLen)

	if offset+2 > len(data) {
		return key, value, 0, fmt.Errorf("incomplete error code in response")
	}

	errCode = int16(binary.BigEndian.Uint16(data[offset : offset+2]))

	return key, value, errCode, nil
}
