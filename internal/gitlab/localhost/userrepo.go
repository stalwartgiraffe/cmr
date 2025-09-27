package localhost

import "fmt"

// InMemoryUserRepository  In-memory implementation
type InMemoryUserRepository struct {
	users  map[int]*User
	nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) Create(user User) (*User, error) {
	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = &user
	return &user, nil
}

func (r *InMemoryUserRepository) GetByID(id int) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func (r *InMemoryUserRepository) List() ([]*User, error) {
	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
