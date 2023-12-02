package masterblaster

import "testing"
import "github.com/aws/aws-sdk-go/service/sts"
import "errors"

func TestSuccessGet(t *testing.T) {
	var acctChan chan workerMessage = make(chan workerMessage, 2)
	var doneChan chan workerMessage = make(chan workerMessage, 2)
	var acctnumber string = "999999999999"
	credsReader := func(creds *sts.Credentials, s string) (*sts.Credentials, error) {
		return &sts.Credentials{}, nil
	}
	plugin := func(creds *sts.Credentials, s string, data APIPlugin) error {
		f, _ := data.PluginContextGetValue("foo")
		if f != "bar" {
			return errors.New("TestSuccessfulGet")
		}
		return nil
	}
	acctChan <- workerMessage{
		acctNumber: acctnumber,
		status:     WAITING,
		detail:     "",
	}

	mb := New(new(string), new(string), new(string), 1)
	mb.PluginContextSetKeyValue("foo", "bar")
	close(acctChan)
	worker(nil, plugin, credsReader, acctChan, doneChan, mb.(APIPlugin))

	msg, ok := <-doneChan

	if ok == false {
		t.Error("worker channel fetch failed")
	}
	if msg.status != SUCCESS {
		t.Error("worker status should have been SUCCESS")
	}
}
func TestMissingKey1(t *testing.T) {
	var acctChan chan workerMessage = make(chan workerMessage, 2)
	var doneChan chan workerMessage = make(chan workerMessage, 2)
	var acctnumber string = "999999999999"
	credsReader := func(creds *sts.Credentials, s string) (*sts.Credentials, error) {
		return &sts.Credentials{}, nil
	}
	plugin := func(creds *sts.Credentials, s string, data APIPlugin) error {
		_, b := data.PluginContextGetValue("foo")
		if b {
			return errors.New("TestMissingKey")
		}
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
		t.Error("worker channel fetch failed")
	}
	if msg.status != SUCCESS {
		t.Error("worker status should have been FAILURE")
	}
}
func TestMissingKey2(t *testing.T) {
	var acctChan chan workerMessage = make(chan workerMessage, 2)
	var doneChan chan workerMessage = make(chan workerMessage, 2)
	var acctnumber string = "999999999999"
	credsReader := func(creds *sts.Credentials, s string) (*sts.Credentials, error) {
		return &sts.Credentials{}, nil
	}
	plugin := func(creds *sts.Credentials, s string, data APIPlugin) error {
		_, b := data.PluginContextGetValue("foo1")
		if b {
			return errors.New("TestMissingKey")
		}
		return nil
	}
	acctChan <- workerMessage{
		acctNumber: acctnumber,
		status:     WAITING,
		detail:     "",
	}

	mb := New(new(string), new(string), new(string), 1)
	mb.PluginContextSetKeyValue("foo", "bar")
	close(acctChan)
	worker(nil, plugin, credsReader, acctChan, doneChan, mb.(APIPlugin))

	msg, ok := <-doneChan

	if ok == false {
		t.Error("worker channel fetch failed")
	}
	if msg.status != SUCCESS {
		t.Error("worker status should have been FAILURE")
	}
}
