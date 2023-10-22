package auth

type UseCase interface {
	GetJWTUser(jwtToken string) (*ContextUser, error)
	GetContextUserKey() ContextUserKey
}
