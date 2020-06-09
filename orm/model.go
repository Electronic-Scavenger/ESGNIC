package orm

type User struct {
	ID       int
	QQ       string
	Name     string // 用户名
	Phone    string
	Email    string
	Password string
	ASN      int
}

type Network struct {
	UserID  int
	User    User
	Network string // 192.168.0.0/24
}
