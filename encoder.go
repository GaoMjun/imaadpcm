package imaadpcm

import (
	"bytes"
	"encoding/binary"
)

type Encoder struct {
	prev_sample int16
	samplesCh   chan []byte
	buffer      []byte
}

func NewEncoder() (encoder *Encoder) {
	encoder = &Encoder{}
	encoder.samplesCh = make(chan []byte)
	return
}

func (self *Encoder) Write(p []byte) (n int, err error) {
	var (
		buffer     = bytes.NewBuffer(p)
		index      int
		cur_sample int16
		delta      int16
		sb         int
		samples    []int
		code       int
		bs         []byte
	)

	for {
		if err = binary.Read(buffer, binary.LittleEndian, &cur_sample); err != nil {
			break
		}
		delta = cur_sample - self.prev_sample

		if delta < 0 {
			delta = -delta
			sb = 8
		} else {
			sb = 0
		}

		code = 4 * int(delta) / step_table[index]
		if code > 7 {
			code = 7
		}

		index += index_adjust[code]
		if index < 0 {
			index = 0
		} else if index > 88 {
			index = 88
		}

		self.prev_sample = cur_sample

		samples = append(samples, code|sb)
	}

	for i := 0; i < len(samples); i += 2 {
		bs = append(bs, byte(samples[i]&0xF<<4)|byte(samples[i+1]&0xF))
	}

	self.samplesCh <- bs
	return
}

func (self *Encoder) Read(p []byte) (n int, err error) {
COPY:
	if len(self.buffer) > 0 {
		if len(self.buffer) <= len(p) {
			n = copy(p, self.buffer)
			self.buffer = nil
			return
		}

		n = copy(p, self.buffer)
		self.buffer = self.buffer[n:]
		return
	}

	self.buffer = <-self.samplesCh
	goto COPY
}

func Encode(p []byte) (o []byte) {
	var (
		buffer      = bytes.NewBuffer(p)
		index       int
		cur_sample  int16
		prev_sample int16
		delta       int16
		sb          int
		samples     []int
		code        int
	)

	for {
		if err := binary.Read(buffer, binary.LittleEndian, &cur_sample); err != nil {
			break
		}
		delta = cur_sample - prev_sample

		if delta < 0 {
			delta = -delta
			sb = 8
		} else {
			sb = 0
		}

		code = 4 * int(delta) / step_table[index]
		if code > 7 {
			code = 7
		}

		index += index_adjust[code]
		if index < 0 {
			index = 0
		} else if index > 88 {
			index = 88
		}

		prev_sample = cur_sample

		samples = append(samples, code|sb)
	}

	for i := 0; i < len(samples); i += 2 {
		o = append(o, byte(samples[i]&0xF<<4)|byte(samples[i+1]&0xF))
	}

	return
}
