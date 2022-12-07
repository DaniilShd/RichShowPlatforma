package repository

type DatabaseRepo interface {
	Authenticate(login, testPassword string) (int, int, string, error)
}
