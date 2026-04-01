package auth

type Principal struct {
	ToolGroup string
	Token     string
}

func Authenticate(_ string) (*Principal, error) {
	return nil, nil
}
