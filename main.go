package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/thomaswhitcomb/masterblaster/masterblaster"
	//"github.com/thomaswhitcomb/masterblaster.git/masterblaster"
)

//import "masterblaster"

// These are used for command line options
var useridPtr *string
var profilePtr *string
var mfaPtr *string
var controlPlaneAcctNumberPtr *string
var workersPtr *int

// Used as a goroutine to get results safely (synchronized) from the ResultPipe.
func resultPipeReader(mb masterblaster.API) {
	for {
		message, ok := mb.ResultPipeGet()
		if ok {
			fmt.Printf("result message: %s\n", message)
		} else {
			fmt.Printf("result stopping\n")
			break
		}
	}
}
func main() {

	useridPtr = flag.String(
		"userid",
		"whitcomb",
		"A userid")

	profilePtr = flag.String(
		"profile",
		"pcco",
		"A profile identifier in your ~/.aws/credentials")

	mfaPtr = flag.String(
		"mfa",
		"none",
		"A six digit integer")

	controlPlaneAcctNumberPtr = flag.String(
		"acctNumber",
		"244268218855",
		"The control plane account number")

	workersPtr = flag.Int(
		"workers",
		10,
		"The number of concurrent workers. Default is 10, maximum is 100")

	flag.Parse()

	mb := masterblaster.New(useridPtr, profilePtr, controlPlaneAcctNumberPtr, *workersPtr)

	go resultPipeReader(mb) // Start the ResultPipe reader

	// Run the plugin (see below)
	mb.PluginContextSetKeyValue("greeting", "hello")
	mb.Run(Plugin, *mfaPtr)

}
func Plugin(creds *sts.Credentials, acctNumber string, mb masterblaster.APIPlugin) error {
	s, _ := mb.PluginContextGetValue("greeting")
	mb.ResultPipePut(s)
	return nil
}
