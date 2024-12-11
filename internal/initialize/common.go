package initialize

import (
	"log"

	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/utils/token"
)

func InitCommon() {
	tokenMaker, err := token.NewPasetoMaker(global.Config.Token.SecretKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	global.TokenMaker = tokenMaker
}
