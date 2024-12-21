package function

// Kind is the kind of function.
//
// TODO(karlhepler): This might need to move somewhere else.
type Kind int

const (
	// KindAdd is a function that was added in a commit.
	KindAdd Kind = iota + 1
	// KindMod is a function that was modified in a commit.
	KindMod
	// KindDel is a function that was deleted in a commit.
	KindDel
)
