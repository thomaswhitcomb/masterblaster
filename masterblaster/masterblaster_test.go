package masterblaster

import "testing"

func TestNew(t *testing.T) {
	var profilePtr string = ""
	var controlPlaneAcctNumber string = ""
	var acctNumber string = ""

	mb := New(&acctNumber, &profilePtr, &controlPlaneAcctNumber, 10)
	if mb == nil {
		t.Error("New for masterblaster failed")
	}
}

func TestResultPutAndGet(t *testing.T) {
	var profilePtr string = ""
	var controlPlaneAcctNumber string = ""
	var acctNumber string = ""

	mb := New(&acctNumber, &profilePtr, &controlPlaneAcctNumber, 10)
	for i := 0; i < 10; i++ {
		mb.(APIPlugin).ResultPipePut(i)
	}
	for i := 0; i < 10; i++ {
		msg, _ := mb.ResultPipeGet()
		x := msg.(int)
		if x != i {
			t.Errorf("ResultPipe Put and Get test failed. x = %d and i = %d\n", x, i)
		}
	}
}
func TestQueueAccounts(t *testing.T) {
	accounts := []string{"1", "2", "3", "4", "5"}
	var c chan workerMessage = make(chan workerMessage, len(accounts))
	queueAccounts(accounts, c)

	if len(c) != len(accounts) {
		t.Errorf("TestQueueAccounts failed.  Bad channel length: %d\n", len(c))
	}
	for i := 0; i < len(accounts); i++ {
		elem := <-c
		if elem.status != WAITING {
			t.Error("TestQueueAccounts failed.  Bad status")
		}
	}
}
func TestDequeueResponses(t *testing.T) {
	accounts := []string{"1", "2", "3", "4", "5"}
	var c chan workerMessage = make(chan workerMessage, len(accounts))
	for i := 0; i < len(accounts); i++ {
		c <- workerMessage{}
	}
	dequeueResponses(accounts, c)
	if len(c) != 0 {
		t.Error("TestDequeue failed.  Channel not empty")
	}
}
