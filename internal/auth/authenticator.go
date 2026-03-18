package auth

type Authenticator struct {
	adminPassword string
}

func New(adminPassword string) Authenticator {
	return Authenticator{
		adminPassword: adminPassword,
	}
}

func (a Authenticator) Authenticate(password string) bool {
	return password == a.adminPassword
}
