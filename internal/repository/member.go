package repository

// Member hold metadata about a Membe￿r
type Member struct {
	ID		int64
	Name	string
}

// MemberDatabase provides access to a database od member
type MemberDatabase interface {
	// AddMember saves a given member, assigning it a new ID
	AddMember(member *Member) (id int64, err error)
}
