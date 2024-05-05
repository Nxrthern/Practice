package token

var _ Verifier = &jwtTokenVerifyer{}

const ()

var ()

func NewJwTVerifier() Verifier {
	return &jwtTokenVerifyer{}
}

type jwtTokenVerifyer struct {
}

func (s *jwtTokenVerifyer) Verify(tokenStr string) error {
	return s.Verify(tokenStr)
}
