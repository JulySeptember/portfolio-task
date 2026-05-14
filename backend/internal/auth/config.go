package auth

import (
	"fmt"
	"os"
)

var (
	cognitoRegion     string
	cognitoUserPoolID string
	cognitoClientID   string
	cognitoIssuer     string
)

func init() {

	cognitoRegion = os.Getenv(
		"AWS_REGION",
	)

	cognitoUserPoolID = os.Getenv(
		"COGNITO_USER_POOL_ID",
	)

	cognitoClientID = os.Getenv(
		"COGNITO_APP_CLIENT_ID",
	)

	if cognitoRegion != "" &&
		cognitoUserPoolID != "" {

		cognitoIssuer = fmt.Sprintf(
			"https://cognito-idp.%s.amazonaws.com/%s",
			cognitoRegion,
			cognitoUserPoolID,
		)
	}
}
