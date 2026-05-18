package env

const (
	ENV_OS_PWD = "PWD"
)

func GetPwd() string {
	return Get(ENV_OS_PWD)
}
