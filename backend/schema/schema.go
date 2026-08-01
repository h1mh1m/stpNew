package schemas

type loginSchema struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type signupSchema struct {
	nama       string ``
	email      string ``
	npwp       string `default:"NO"`
	password   string ``
	isVerified bool   `default:"false"`
}
