package valorant

import "errors"

var (
	ErrorRiotAuthentication = errors.New("riot_authentication_error")
	ErrorRiotMultifactor    = errors.New("riot_multifactor_error")
	ErrorRiotRateLimit      = errors.New("riot_ratelimit_error")

	ErrorRiotUnknownResponseType = errors.New("riot_unknown_response_type_error")
	ErrorRiotUnknownErrorType    = errors.New("riot_unknown_error_type_error")

	ResponseErrors = map[string]error{
		"auth_failure": ErrorRiotAuthentication,
		"rate_limited": ErrorRiotRateLimit,
	}
)
