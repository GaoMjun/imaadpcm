package imaadpcm

type Decoder struct {
	samplesCh chan []byte
	buffer    []byte
}

func NewDecoder() (decoder *Decoder) {
	decoder = &Decoder{}
	decoder.samplesCh = make(chan []byte)
	return
}

func (self *Decoder) Write(p []byte) (n int, err error) {
	var (
		index int
		code  int
		bs    []byte
	)

	for _, v := range p {
		code1 := int(v & 0xF0 >> 4)
		code2 := int(v & 0x0F)

		code, index = decode(code1, index)
		bs = append(bs, byte(code&0x00FF))
		bs = append(bs, byte(code&0xFF00>>8))

		code, index = decode(code2, index)
		bs = append(bs, byte(code&0x00FF))
		bs = append(bs, byte(code&0xFF00>>8))
	}

	self.samplesCh <- bs
	return
}

func (self *Decoder) Read(p []byte) (n int, err error) {
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

func decode(code_i, index_i int) (code_o, index_o int) {
	var (
		sb    int
		delta int
	)

	if code_i&8 != 0 {
		sb = 1
	} else {
		sb = 0
	}

	code_i &= 7

	delta = (step_table[index_i]*code_i)/4 + step_table[index_i]/8
	if sb == 1 {
		delta = -delta
	}

	code_o += delta
	if code_o > 32767 {
		code_o = 32767
	} else if code_o < -32768 {
		code_o = -32768
	}

	index_i += index_adjust[code_i]
	if index_i < 0 {
		index_i = 0
	} else if index_i > 88 {
		index_i = 88
	}

	index_o = index_i
	return
}

func Decode(p []byte) (o []byte) {
	var (
		index int
		code  int
	)

	for _, v := range p {
		code1 := int(v & 0xF0 >> 4)
		code2 := int(v & 0x0F)

		code, index = decode(code1, index)
		o = append(o, byte(code&0x00FF))
		o = append(o, byte(code&0xFF00>>8))

		code, index = decode(code2, index)
		o = append(o, byte(code&0x00FF))
		o = append(o, byte(code&0xFF00>>8))
	}

	return
}
