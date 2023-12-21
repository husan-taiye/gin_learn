package user

const (
	emailRegexPattern    = `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	nicknameRegexPattern = `^.{1,64}$`
	birthdayRegexPattern = `^\d{4}-\d{2}-\d{2}$`
	profileRegexPattern  = `^.{0,255}$`
)
