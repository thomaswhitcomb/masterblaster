package masterblaster

import "testing"
import "github.com/aws/aws-sdk-go/service/sts"
import "errors"
import "strings"

func TestFailingCustomerCredentialReader(t *testing.T) {
	var acctChan chan workerMessage = make(chan workerMessage, 2)
	var doneChan chan workerMessage = make(chan workerMessage, 2)
	var acctnumber string = "999999999999"
	credsReader := func(creds *sts.Credentials, s string) (*sts.Credentials, error) {
		return nil, errors.New("TestFailingCustomerCredentialsReader")
	}
	acctChan <- workerMessage{
		acctNumber: acctnumber,
		status:     WAITING,
		detail:     "",
	}

	mb := New(new(string), new(string), new(string), 1)

	close(acctChan)
	worker(nil, nil, credsReader, acctChan, doneChan, mb.(APIPlugin))

	msg, ok := <-doneChan

	if ok == false {
		t.Error("worker channel fetch failed.")
	}
	if msg.status != FAILURE {
		t.Error("worker status should have been FAILURE")
	}

	if strings.HasPrefix(msg.detail, "Failed to assume role into account") == false {
		t.Errorf("worker detail unexpected: %s\n ", msg.detail)
	}

}
func TestFailingPlugin(t *testing.T) {
	var acctChan chan workerMessage = make(chan workerMessage, 2)
	var doneChan chan workerMessage = make(chan workerMessage, 2)
	var acctnumber string = "999999999999"
	credsReader := func(creds *sts.Credentials, s string) (*sts.Credentials, error) {
		return &sts.Credentials{}, nil
	}
	plugin := func(creds *sts.Credentials, s string, data APIPlugin) error {
		return errors.New("TestFailingCustomerCredentialsReader")
	}
	acctChan <- workerMessage{
		acctNumber: acctnumber,
		status:     WAITING,
		detail:     "",
	}

	mb := New(new(string), new(string), new(string), 1)

	close(acctChan)
	worker(nil, plugin, credsReader, acctChan, doneChan, mb.(APIPlugin))

	msg, ok := <-doneChan

	if ok == false {
		t.Error("worker channel fetch failed.")
	}
	if msg.status != FAILURE {
		t.Error("worker status should have been FAILURE")
	}
}
func TestSuccessfulPlugin(t *testing.T) {
	var acctChan chan workerMessage = make(chan workerMessage, 2)
	var doneChan chan workerMessage = make(chan workerMessage, 2)
	var acctnumber string = "999999999999"
	credsReader := func(creds *sts.Credentials, s string) (*sts.Credentials, error) {
		return &sts.Credentials{}, nil
	}
	plugin := func(creds *sts.Credentials, s string, data APIPlugin) error {
		return nil
	}
	acctChan <- workerMessage{
		acctNumber: acctnumber,
		status:     WAITING,
		detail:     "",
	}

	mb := New(new(string), new(string), new(string), 1)

	close(acctChan)
	worker(nil, plugin, credsReader, acctChan, doneChan, mb.(APIPlugin))

	msg, ok := <-doneChan

	if ok == false {
		t.Error("worker channel fetch failed.")
	}
	if msg.status != SUCCESS {
		t.Error("worker status should have been FAILURE")
	}
	if strings.HasPrefix(msg.detail, "Completed successfully") == false {
		t.Errorf("worker detail unexpected: %s\n ", msg.detail)
	}
}
