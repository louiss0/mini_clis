package task

type Task struct {
	Title,
	Description,
	UpdatedAt,
	createdAt string
}

func (self Task) CreatedAt() string {

	return self.createdAt
}
