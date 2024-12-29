package github

type Owner struct {
	Login string
}

func (o Owner) String() string {
	return o.Login
}
