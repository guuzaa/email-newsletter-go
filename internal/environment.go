package internal

type Environment string

const (
	Local      Environment = "local"
	Production Environment = "production"
)

func (e Environment) String() string {
	return string(e)
}

func ParseEnvironment(str string) Environment {
	switch str {
	case "local":
		return Local
	case "production":
		return Production
	default:
		return Local
	}
}
