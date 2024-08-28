package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_USER     = "root"
	DB_PASSWORD = ""
	DB_NAME     = "mini_challenge_go_6"
)

func connectDB() *sql.DB {
	mysqlInfo := fmt.Sprintf("%s:%s@/%s", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("mysql", mysqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type Product struct {
	ID        int
	Name      string
	CreatedAt string
	UpdatedAt string
}

type Variant struct {
	ID          int
	VariantName string
	Quantity    int
	CreatedAt   string
	UpdatedAt   string
}

type ProductWithVariants struct {
	ID        int
	Name      string
	CreatedAt string
	UpdatedAt string
	Variants  []Variant
}

// Create Product and get auto-incremented ID
func createProduct(db *sql.DB, name string) (int64, error) {
	query := "INSERT INTO products (name, created_at, updated_at) VALUES (?, NOW(), NOW())"
	result, err := db.Exec(query, name)
	if err != nil {
		return 0, err
	}

	// Get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func updateProduct(db *sql.DB, id int, name string) error {
	query := "UPDATE products SET name=?, updated_at=NOW() WHERE id=?"
	_, err := db.Exec(query, name, id)
	return err
}

func getProductById(db *sql.DB, id int) (*Product, error) {
	var product Product
	query := "SELECT id, name, created_at, updated_at FROM products WHERE id=?"
	row := db.QueryRow(query, id)
	err := row.Scan(&product.ID, &product.Name, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create Variant and get auto-incremented ID
func createVariant(db *sql.DB, variantName string, quantity int, productId int) (int64, error) {
	query := "INSERT INTO variants (variant_name, quantity, product_id, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
	result, err := db.Exec(query, variantName, quantity, productId)
	if err != nil {
		return 0, err
	}

	// Get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func updateVariantById(db *sql.DB, id int, variantName string, quantity int) error {
	query := "UPDATE variants SET variant_name=?, quantity=?, updated_at=NOW() WHERE id=?"
	_, err := db.Exec(query, variantName, quantity, id)
	return err
}

func deleteVariantById(db *sql.DB, id int) error {
	query := "DELETE FROM variants WHERE id=?"
	_, err := db.Exec(query, id)
	return err
}

func getProductWithVariants(db *sql.DB, productId int) (*ProductWithVariants, error) {
	var product ProductWithVariants

	query := "SELECT id, name, created_at, updated_at FROM products WHERE id=?"
	row := db.QueryRow(query, productId)
	err := row.Scan(&product.ID, &product.Name, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}

	query = "SELECT id, variant_name, quantity, created_at, updated_at FROM variants WHERE product_id=?"
	rows, err := db.Query(query, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var variant Variant
		err := rows.Scan(&variant.ID, &variant.VariantName, &variant.Quantity, &variant.CreatedAt, &variant.UpdatedAt)
		if err != nil {
			return nil, err
		}
		product.Variants = append(product.Variants, variant)
	}

	return &product, nil
}

func main() {
	db := connectDB()
	defer db.Close()

	productID, err := createProduct(db, "Sample Product")
	if err != nil {
		log.Fatal(err)
	}

	variantID, err := createVariant(db, "Sample Variant", 100, int(productID))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created Variant with ID: %d\n", variantID)

	productWithVariants, err := getProductWithVariants(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product with Variants: %+v\n", productWithVariants)
}
