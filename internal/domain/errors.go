package domain

import "fmt"

var (
	// Internal database error
	ErrInternalDatabase = fmt.Errorf("internal database error")
	ErrSQLNoRows        = fmt.Errorf("sql no rows")

	// No user in db
	ErrNoUser = fmt.Errorf("no user error")

	// No client_id
	ErrNoClientID = fmt.Errorf("no client_id error")

	// No fcm_token
	ErrNoFCMToken = fmt.Errorf("no fcm_token error")
	// Incorrect os
	ErrIncorrectOs = fmt.Errorf("os != android and os != ios")
	// PushFailed
	ErrPushFailed = fmt.Errorf("push failed")

	// Unauthorized
	ErrUnauthorized = fmt.Errorf("unauthorized")

	// Bad request
	ErrInvalidInputData       = fmt.Errorf("invalid input data")
	ErrPaymentsLessThanAmount = fmt.Errorf("payments less than amount")
	ErrTooHighPayment         = fmt.Errorf("too high payment")
	ErrAmountLessThanPayment  = fmt.Errorf("amount less than payment")

	// Validation Failed
	ErrValidationFailed = fmt.Errorf("validation failed")

	ErrUnconfirmedEmail = fmt.Errorf("unconfirmed email")
	ErrInvalidSMSCode   = fmt.Errorf("invalid sms code")
	ErrDuplicateRequest = fmt.Errorf("duplicate request")
	ErrSMSSending       = fmt.Errorf("sms sending failed")

	// Internal Security error
	ErrInternalSecurity           = fmt.Errorf("internal security error")
	ErrSecurityUnsupportedKeyType = fmt.Errorf("unsupported key type")

	// Internal integration error
	ErrInternalIntegration = fmt.Errorf("internal integration error")

	// Internal push-service error
	ErrInternalPushService = fmt.Errorf("internal push service error")

	// Internal OpenAPI error
	ErrInternalOpenAPI = fmt.Errorf("internal openapi error")

	ErrUnavailableBureauReport = fmt.Errorf("this bureau report is currenty unavailable")

	ErrNotFound = fmt.Errorf("not found")
)
