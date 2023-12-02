package masterblaster

type resultPipe interface {
	put(data interface{})
	get() (interface{}, bool)
	len() int
	close()
}
type pipe struct {
	p chan interface{}
}

func newResultPipe(buffersize int) resultPipe {
	s := &pipe{
		p: make(chan interface{}, buffersize),
	}
	return s
}

func (s *pipe) close() {
	close(s.p)
}

func (s *pipe) len() int {
	return len(s.p)
}

func (s *pipe) put(data interface{}) {
	s.p <- data
}

func (s *pipe) get() (interface{}, bool) {
	msg, ok := <-s.p
	return msg, ok
}
