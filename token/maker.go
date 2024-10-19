package token

import "time"

type Maker interface {
	// CreateToken creates a new token for sepcific and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the tokern is valid or not 
	VerifyToken (token string) (*Payload, error)
}