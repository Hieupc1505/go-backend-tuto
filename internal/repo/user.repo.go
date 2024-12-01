package repos



type Repo struct{}

func UserRepo() *Repo{
	return &Repo{}; 
}


func (*Repo) GetUserInfo() string{
	return "Tipsgo"
}