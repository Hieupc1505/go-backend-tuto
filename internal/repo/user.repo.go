package repos

import (
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
)

// type Repo struct{}

// func UserRepo() *Repo{
// 	return &Repo{};
// }

// func (*Repo) GetUserInfo() string{
// 	return "Tipsgo"
// }

type IUserRepository interface {
	GetUserByEmail(email string) (db.User, error)
}

type userRepository struct{}

func (*userRepository) GetUserByEmail(email string) (db.User, error) {
	user, err := global.PgDb.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func NewUserRepository() IUserRepository {
	return &userRepository{}
}
