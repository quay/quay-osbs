package googlelogin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/quay/config-tool/pkg/lib/shared"
)

// Validate checks the configuration settings for this field group
func (fg *GoogleLoginFieldGroup) Validate(opts shared.Options) []shared.ValidationError {

	fgName := "GoogleLogin"

	var errors []shared.ValidationError

	// If google login is off, return false
	if fg.FeatureGoogleLogin == false {
		return errors
	}

	// Check for config
	if fg.GoogleLoginConfig == nil {
		newError := shared.ValidationError{
			Tags:       []string{"GOOGLE_LOGIN_CONFIG"},
			FieldGroup: fgName,
			Message:    "GOOGLE_LOGIN_CONFIG is required for GoogleLogin",
		}
		errors = append(errors, newError)
		return errors
	}

	// Check for client id
	if fg.GoogleLoginConfig.ClientId == "" {
		newError := shared.ValidationError{
			Tags:       []string{"GOOGLE_LOGIN_CONFIG.CLIENT_ID"},
			FieldGroup: fgName,
			Message:    "GOOGLE_LOGIN_CONFIG.CLIENT_ID is required for GoogleLogin",
		}
		errors = append(errors, newError)
	}

	// Check for endpoint
	if fg.GoogleLoginConfig.ClientSecret == "" {
		newError := shared.ValidationError{
			Tags:       []string{"GOOGLE_LOGIN_CONFIG.CLIENT_SECRET"},
			FieldGroup: fgName,
			Message:    "GOOGLE_LOGIN_CONFIG.CLIENT_SECRET is required for GoogleLogin",
		}
		errors = append(errors, newError)
	}

	// Check OAuth endpoint
	success := ValidateGoogleOAuth(fg.GoogleLoginConfig.ClientId, fg.GoogleLoginConfig.ClientSecret)
	if !success {
		newError := shared.ValidationError{
			Tags:       []string{"GOOGLE_LOGIN_CONFIG.CLIENT_ID", "GOOGLE_LOGIN_CONFIG.CLIENT_SECRET"},
			FieldGroup: fgName,
			Message:    "Could not verify Google OAuth credentials",
		}
		errors = append(errors, newError)
	}

	return errors
}

// ValidateGoogleOAuth checks that the Bitbucker OAuth credentials are correct
func ValidateGoogleOAuth(clientID, clientSecret string) bool {

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("POST", "https://www.googleapis.com/oauth2/v3/token?client_id="+clientID+"&client_secret="+clientSecret+"&grant_type=authorization_code&code=FAKECODE&redirect_uri=https://fakeredirect.com", nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	// Load response into json
	var responseJSON map[string]interface{}
	_ = json.Unmarshal(respBody, &responseJSON)

	// If the error isnt unauthorized
	if responseJSON["error"] == "invalid_grant" {
		return true
	}

	return false

}
