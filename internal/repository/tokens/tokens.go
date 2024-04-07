package tokens

type Storage interface {
	Disconnect() error
	WriteRefreshToken(id int, refreshToken string) error
	UpdateRefreshToken(id int, refreshToken string) error
	ReadRefreshToken(id int) (string, error)
}
