package user

const (
	emailRegexPattern    = `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
)
