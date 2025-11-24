package user

import "fmt"

type UserModule struct{}

func New() *UserModule { return &UserModule{} }

func (m *UserModule) Name() string { return "user" }
func (m *UserModule) Init()        { fmt.Println("User module initialized") }
