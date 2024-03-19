package main

import (
	"log"

	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/account"
	"github.com/ednaldo-dilorenzo/iappointment/modules/auth"
	"github.com/ednaldo-dilorenzo/iappointment/modules/category"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/modules/transaction"
	"github.com/ednaldo-dilorenzo/iappointment/server"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(auth.NewAuthRepository)
	container.Provide(auth.NewAuthService)
	container.Provide(auth.NewAuthController)
	container.Provide(generic.NewGenericRepository[model.Category])
	container.Provide(generic.NewGenericService[model.Category])
	container.Provide(category.NewCategoryController)
	container.Provide(generic.NewGenericRepository[model.Account])
	container.Provide(generic.NewGenericService[model.Account])
	container.Provide(account.NewAccountController)
	container.Provide(generic.NewGenericRepository[model.Transaction])
	container.Provide(generic.NewGenericService[model.Transaction])
	container.Provide(transaction.NewTransactionController)
	container.Provide(server.NewServer)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(server *server.Server) error {
		waitforShutdownInterrupt := server.Start()
		<-waitforShutdownInterrupt

		log.Println("Shutting Down Server..")

		server.ShutdownGracefully()

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
