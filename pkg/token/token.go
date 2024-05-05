package token

type Verifier interface {
	Verify(str string, modUser string) error
}
