package token

type Verifier interface {
	Verify(token string) error
}
