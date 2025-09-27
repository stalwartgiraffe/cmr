package localhost

// Service layer with dependency injection
type UserService struct {
	repo UserRepository
}

// UserRepository Repository interface for dependency injection
type UserRepository interface {
	Create(user User) (*User, error)
	GetByID(id int) (*User, error)
	List() ([]*User, error)
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}
