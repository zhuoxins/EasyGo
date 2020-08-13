package validate

import "regexp"

var (
	emailPattern  = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)
	ipPattern     = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	mobilePattern = regexp.MustCompile(`^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][01356789]|[4][579]))\d{8}$`)
	telPattern    = regexp.MustCompile(`^(0\d{2,3}(\-)?)?\d{7,8}$`)
)

type verifyFunc func(string, string, map[string]string) (bool, string)
