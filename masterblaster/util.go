package masterblaster

import "fmt"
import "log"
import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/aws/credentials"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/service/sts"

// Assumes into the PCC control plane and returns temp credentials
func getControlPlaneCredentials(
	profile string,
	acctNumber string,
	serialNumber string,
	tokenCode string) (*sts.Credentials, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-west-2")},
		Profile: profile,
	}))
	svc := sts.New(sess)
	arn := fmt.Sprintf("arn:aws:iam::%s:role/Manage/pcccp-cloudenabler-human-super-admin", acctNumber)
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),                                                       // Required
		RoleSessionName: aws.String(fmt.Sprintf("%s-pcc-remove-extra-cloudtrail", acctNumber)), // Required
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(tokenCode),
	}
	resp, err := svc.AssumeRole(params)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	return resp.Credentials, nil
}

// Assume into the customer's account and return temp credentials.
func getCustomerCredentials(
	creds *sts.Credentials,
	acctNumber string) (*sts.Credentials, error) {

	var sess *session.Session = session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken),
	}))
	var svc *sts.STS = sts.New(sess)
	var arn string = fmt.Sprintf("arn:aws:iam::%s:role/PCC/pcc-cloudenabler-human-admin", acctNumber)
	var params *sts.AssumeRoleInput = &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn), // Required
		RoleSessionName: aws.String(fmt.Sprintf("%s-cloudtrail-update", acctNumber)),
	}
	resp, err := svc.AssumeRole(params)
	if err != nil {
		return nil, err
	}
	return resp.Credentials, nil
}
