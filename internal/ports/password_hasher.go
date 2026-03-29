package ports

type IPasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}