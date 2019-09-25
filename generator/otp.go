package generator

// OTP ...
func OTP() string {
	num, errRand := RandomNumericString(6)
	if errRand != nil {
		return ""
	}
	return num
}
