package masterblaster

import "testing"

func TestNewResultPipe(t *testing.T) {
	p := newResultPipe(10)
	if p == nil {
		t.Error("ResultPipe creation failed")
	}
}
func TestBuffering(t *testing.T) {
	pipe := newResultPipe(10)
	for i := 0; i < 10; i++ {
		pipe.put(i)
	}
	if pipe.len() != 10 {
		t.Error("ResultPipe buffer test failed")
	}
}
func TestPutAndGet(t *testing.T) {
	pipe := newResultPipe(10)
	for i := 0; i < 10; i++ {
		pipe.put(i)
	}
	for i := 0; i < 10; i++ {
		msg, _ := pipe.get()
		x := msg.(int)
		if x != i {
			t.Errorf("ResultPipe Put and Get test failed. x = %d and i = %d\n", x, i)
		}
	}
}
