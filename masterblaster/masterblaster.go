package masterblaster

import "fmt"

import "log"
import "bufio"
import "os"
import "github.com/aws/aws-sdk-go/service/sts"

type API interface {
	ResultPipeGet() (interface{}, bool)
	PluginContextSetKeyValue(key string, value string)
	Run(plugin Plugin, mfa string)
}
type APIPlugin interface {
	ResultPipePut(s interface{})
	PluginContextGetValue(key string) (string, bool)
}

type Plugin func(creds *sts.Credentials, acctNumber string, data APIPlugin) error

type pluginContext map[string]string
type customerCredentialReader func(creds *sts.Credentials, acctNumber string) (*sts.Credentials, error)

type masterBlaster struct {
	userid                 string
	profile                string
	mfa                    string
	controlPlaneAcctNumber string
	workers                int
	pipe                   resultPipe
	context                pluginContext
}

func New(useridPtr *string, profilePtr *string, controlPlaneAcctNumberPtr *string, workers int) API {
	mb := &masterBlaster{
		userid:                 *useridPtr,
		profile:                *profilePtr,
		controlPlaneAcctNumber: *controlPlaneAcctNumberPtr,
		workers:                workers,
		pipe:                   newResultPipe(workers),
		context:                make(pluginContext),
	}
	return mb
}

func (mb *masterBlaster) ResultPipePut(s interface{}) {
	mb.pipe.put(s)
}
func (mb *masterBlaster) ResultPipeGet() (interface{}, bool) {
	return mb.pipe.get()
}
func (mb *masterBlaster) PluginContextGetValue(key string) (string, bool) {
	s, b := mb.context[key]
	return s, b
}
func (mb *masterBlaster) PluginContextSetKeyValue(key string, value string) {
	mb.context[key] = value
}
func (mb *masterBlaster) Run(plugin Plugin, mfa string) {

	defer mb.pipe.close()

	// Load up the accounts.
	var accts []string = loadUpAccts()

	// Get control plane credentials
	creds, err := getControlPlaneCredentials(
		mb.profile,
		mb.controlPlaneAcctNumber,
		fmt.Sprintf(
			"arn:aws:iam::%s:mfa/%s",
			mb.controlPlaneAcctNumber,
			mb.userid), mfa)

	if err != nil {
		log.Fatalf("Control plane failure: %v", err)
	}

	// Input channel for the workers.  Filled with account numbers
	var acctChan chan workerMessage = make(chan workerMessage, len(accts))

	// Output channel from the workers.  Has status for each account.
	var doneChan chan workerMessage = make(chan workerMessage, len(accts))

	// Queue up the account numbers for processing
	queueAccounts(accts, acctChan)

	fn := func() {
		worker(creds, plugin, getCustomerCredentials, acctChan, doneChan, mb)
	}

	startWorkers(mb.workers, fn)

	dequeueResponses(accts, doneChan)
}

// Reads accounts from standard input.  One per line.
func loadUpAccts() []string {
	var accts []string = make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		acct := scanner.Text()
		accts = append(accts, acct)
	}
	return accts
}

func queueAccounts(accounts []string, acctChan chan workerMessage) {
	// Queue up the account numbers for processing
	for i := 0; i < len(accounts); i++ {
		log.Printf("queuing '%s'\n", accounts[i])
		m := workerMessage{
			acctNumber: accounts[i],
			status:     WAITING,
			detail:     "",
		}
		acctChan <- m
	}
}

func startWorkers(count int, fn func()) {
	for i := 0; i < count; i++ {
		go fn()
	}
}

func dequeueResponses(accounts []string, doneChan chan workerMessage) {
	// Receive all the responses from the workers
	for i := 0; i < len(accounts); i++ {
		cc := <-doneChan
		log.Println(cc)
	}
}
