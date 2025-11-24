package user

import (
	"neonexcore/internal/config"
	"neonexcore/internal/core"
	"neonexcore/pkg/database"
)

func (m *UserModule) RegisterServices(c *core.Container) {
	// Register Repository
	c.Provide(func() *UserRepository {
		db := config.DB.GetDB()
		return NewUserRepository(db)
	}, core.Singleton)

	// Register Transaction Manager
	c.Provide(func() *database.TxManager {
		db := config.DB.GetDB()
		return database.NewTxManager(db)
	}, core.Singleton)

	// Register Service
	c.Provide(func() *UserService {
		repo := core.Resolve[*UserRepository](c)
		txManager := core.Resolve[*database.TxManager](c)
		return NewUserService(repo, txManager)
	}, core.Singleton)

	// Register Controller
	c.Provide(func() *UserController {
		service := core.Resolve[*UserService](c)
		return NewUserController(service)
	}, core.Transient)
}
