package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"zocket/image_processing"
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
	
	var ProductBody struct{
		UserID int `json:"user_id"`
		ProductName        string   `json:"product_name"`
		ProductDescription string   `json:"product_description"`
		ProductImages      models.StringArray   `json:"product_images" gorm:"type:json"`
		ProductPrice       float64  `json:"product_price"`
	}


	err := context.BodyParser(&ProductBody)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	var localFilePaths, compressedFilePaths models.StringArray
    for _, imageUrl := range ProductBody.ProductImages {
        downloadPath := "downloads/" + filepath.Base(imageUrl)
        compressedPath := "compressed/" + filepath.Base(imageUrl)
        bucketName := "pet-adopt-72ad5.appspot.com"
        objectName := "images/" + filepath.Base(imageUrl)

        uploadedURL, err := image_processing.ProcessImage(imageUrl, downloadPath, compressedPath, bucketName, objectName)
        if err != nil {
            log.Fatal("Error:", err.Error())
        }

        localFilePaths = append(localFilePaths, imageUrl)
        compressedFilePaths = append(compressedFilePaths, uploadedURL)
    }

	finalProduct:=models.Products{
		UserID:ProductBody.UserID,
		ProductName:ProductBody.ProductName,
		ProductDescription:ProductBody.ProductDescription,
		ProductImages:localFilePaths,
		ProductPrice:ProductBody.ProductPrice,
		CompressedImages:compressedFilePaths,
	}

	

	err = r.DB.Create(&finalProduct).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create product"})
		return err
	}
	
	

	if err != nil {
		log.Fatal("Error:", err.Error())
	}


	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "product has been added",
		"data":    finalProduct,})
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