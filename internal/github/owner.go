package github

type Owner string

func (o Owner) String() string {
	return string(o)
}
