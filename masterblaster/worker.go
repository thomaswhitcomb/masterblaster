package masterblaster

import "log"
import "fmt"
import "github.com/aws/aws-sdk-go/service/sts"

const (
	FAILURE = -1
	WAITING = 0
	SUCCESS = 1
)

type workerMessage struct {
	acctNumber string
	status     int // -1 failed, 0 not processed yet, 1 successful
	detail     string
}

// This runs concurrently as multiple instances.  Reads an account
// from a channel, sends it to the plugin and then responds to the output channel
func worker(
	creds *sts.Credentials,
	plugin Plugin,
	credentialReader customerCredentialReader,
	acctChan chan workerMessage,
	doneChan chan workerMessage,
	data APIPlugin) {

	log.Printf("worker started\n")

	for {
		message, ok := <-acctChan
		if ok {
			custCreds, err := credentialReader(creds, message.acctNumber)
			if err != nil {
				message.status = FAILURE
				message.detail = fmt.Sprintf("Failed to assume role into account. %v", err)
			} else {
				// CALL THE PLUGIN
				if err := plugin(custCreds, message.acctNumber, data); err == nil {
					message.status = SUCCESS
					message.detail = "Completed successfully"
				} else {
					message.status = FAILURE
					message.detail = fmt.Sprintf("%v", err)
				}
			}
			doneChan <- message
		} else {
			break
		}
	}
	doneChan <- workerMessage{
		acctNumber: "",
		status:     FAILURE,
		detail:     "worker ternminating due to bad input channel",
	}
}
