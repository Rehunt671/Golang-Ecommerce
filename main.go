package main



import(
	"github.com/kingsglaive/controllers"
	"github.com/kingsglaive/database"
	"github.com/kingsglaive/middleware"
	"github.com/kingsglaive/routes"
	"github.com/gin-gonic/gin"
)


func main(){

	port := os.Getenv("PORT")
	if port == ""{
		port ="8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client,"Products"),database.UserData(database.Client,"Users"))
	
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.Get("/addtocart",app.AddtoCart())
	router.Get("/removeitem",app.RemoveItem())
	router.Get("/cartcheckout",app.BuyFromCart())
	router.Get("/instantbuy",app.Instantbuy())

	log.Fatal(router.Run(":" + port))


}

