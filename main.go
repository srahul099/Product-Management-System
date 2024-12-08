package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"zocket/models"
	"zocket/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)



type Repository struct{
	DB *gorm.DB
}

func (r *Repository) CreateProduct(context *fiber.Ctx) error {
	product := models.Products{}

	err := context.BodyParser(&product)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&product).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create product"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "product has been added"})
	return nil
}

func (r *Repository) CreateUser(context *fiber.Ctx) error {
	user := models.Users{}

	err := context.BodyParser(&user)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create user"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user has been added"})
	return nil
}



func (r *Repository) GetAllProducts(context *fiber.Ctx) error {
	productModels := &[]models.Products{}

	err := r.DB.Find(productModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get products"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "products fetched successfully",
		"data":    productModels,
	})
	return nil
}

func (r *Repository) GetProducts(context *fiber.Ctx) error {

	id := context.Params("id")
	productModel := &models.Products{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(productModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "product id fetched successfully",
		"data":    productModel,
	})
	return nil
}

func(r *Repository) SetUpRoutes(app *fiber.App){
	api:= app.Group("/api")
	api.Post("/create_user", r.CreateUser)
	api.Post("/create_products", r.CreateProduct)
	api.Get("/products/:id", r.GetProducts)
	api.Get("/products", r.GetAllProducts)
}

func main(){
	err:= godotenv.Load()
	if err!=nil{
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db,err := storage.NewConnection(config)
	if err !=nil{
		log.Fatal("Error connecting to database")
	}


	r := Repository{
		DB:db,
	}

	err = models.MigrateProducts(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	app := fiber.New()
	r.SetUpRoutes(app)
	app.Listen(":3000")
}