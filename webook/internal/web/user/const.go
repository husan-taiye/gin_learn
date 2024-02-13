package user

const (
	EmailRegexPattern    = `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	PasswordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	NicknameRegexPattern = `^.{1,64}$`
	BirthdayRegexPattern = `^\d{4}-\d{2}-\d{2}$`
	ProfileRegexPattern  = `^.{0,255}$`
)
