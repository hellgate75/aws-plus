package access

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"os/user"
)

func HomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		return os.TempDir()
	}
	return usr.HomeDir
}

func AwsCredentialsFilePath() string {
	return fmt.Sprintf("%s%c%s%c%s", HomeFolder(), os.PathSeparator, ".aws", os.PathSeparator, "credentials")
}

func GetOsEnvCredentials() *credentials.Credentials {
	return credentials.NewEnvCredentials()
}

func GetOsCredentials() *credentials.Credentials {
	return credentials.NewSharedCredentials(AwsCredentialsFilePath(), "default")
}

func ReadValueFromCredentials(creds *credentials.Credentials) (*credentials.Value, error) {
	if creds == nil {
		return nil, errors.New(fmt.Sprint("Invalid Credentials object"))
	}
	// Retrieve the credentials value
	val, err := creds.Get()
	return &val, err
}