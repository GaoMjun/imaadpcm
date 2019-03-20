package imaadpcm

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"testing"
)

func TestIMAADPCM(t *testing.T) {
	sample := 2000
	sample2 := -2000

	sampleB := make([]byte, 4)
	sampleB[0] = byte(sample & 0x00FF >> 0)
	sampleB[1] = byte(sample & 0xFF00 >> 8)
	sampleB[2] = byte(sample2 & 0x00FF >> 0)
	sampleB[3] = byte(sample2 & 0xFF00 >> 8)

	encoder := NewEncoder()
	decoder := NewDecoder()

	buffer := make([]byte, 1024)
	n := 0

	go func() {
		n, _ = encoder.Read(buffer)
		decoder.Write(buffer[:n])
	}()

	go func() {
		n, _ = decoder.Read(buffer)
		log.Println(toInt16Array(buffer[:n]))
	}()

	go encoder.Write(sampleB)

	select {}
}

func TestEncodeDecode(t *testing.T) {
	sample := 2000
	sample2 := -2000

	sampleB := make([]byte, 4)
	sampleB[0] = byte(sample & 0x00FF >> 0)
	sampleB[1] = byte(sample & 0xFF00 >> 8)
	sampleB[2] = byte(sample2 & 0x00FF >> 0)
	sampleB[3] = byte(sample2 & 0xFF00 >> 8)

	encodeSample := Encode(sampleB)
	decodeSample := Decode(encodeSample)

	log.Println(toInt16Array(decodeSample))
}

func toInt16Array(data []byte) (o []int16) {
	r := bytes.NewReader(data)
	for {
		var f int16
		if err := binary.Read(r, binary.LittleEndian, &f); err == io.EOF {
			return
		}
		o = append(o, f)
	}

	return
}
