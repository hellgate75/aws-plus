package access

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/hellgate75/aws-cli-samples/access"
)
func ServiceList() {
	svc := servicediscovery.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))
	input := &servicediscovery.ListServicesInput{}

	result, err := svc.ListServices(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case servicediscovery.ErrCodeInvalidInput:
				fmt.Println("Error:", servicediscovery.ErrCodeInvalidInput, awsErr.Error())
			default:
				fmt.Println("Error:", awsErr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println("Error:", err.Error())
		}
		return
	}
	fmt.Println(result)}

func S3ReadAllBuckets() {
	s3svc := s3.New(session.New(&aws.Config{Region: aws.String("eu-west-1"), Credentials: access.GetOsCredentials()}))
	//s3svc.SigningRegion = "eu-west-1"
	//s3svc.Endpoint = "https://s3.eu-west-1.amazonaws.com"
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{
	})
	if err != nil {
		fmt.Println("Failed to list buckets", err)
		return
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	}
}