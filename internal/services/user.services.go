package services

import repos "hieupc05.github/backend-server/internal/repo"

type Services struct {
	info *repos.Repo
}

func NewService() *Services {

	return &Services{
		info: repos.UserRepo(),
	}
}

func (c *Services) GetUserInfoService() string {
	return c.info.GetUserInfo()
}
