package seed

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func Run(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) {
	var count int
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		logger.Error("seed: failed to check users", slog.String("error", err.Error()))
		return
	}
	if count > 0 {
		logger.Info("seed: data already exists, skipping")
		return
	}

	logger.Info("seed: inserting demo data")

	tx, err := pool.Begin(ctx)
	if err != nil {
		logger.Error("seed: failed to begin tx", slog.String("error", err.Error()))
		return
	}
	defer tx.Rollback(ctx)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	hashStr := string(hash)

	var sellerID, buyerID int
	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		"seller@example.com", hashStr, "TechStore Pro", "seller",
	).Scan(&sellerID)
	if err != nil {
		logger.Error("seed: failed to insert seller", slog.String("error", err.Error()))
		return
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		"buyer@example.com", hashStr, "Jane Buyer", "buyer",
	).Scan(&buyerID)
	if err != nil {
		logger.Error("seed: failed to insert buyer", slog.String("error", err.Error()))
		return
	}

	seedCategories(ctx, tx, logger)
	seedProducts(ctx, tx, sellerID, logger)

	if err := tx.Commit(ctx); err != nil {
		logger.Error("seed: failed to commit", slog.String("error", err.Error()))
		return
	}
	logger.Info("seed: demo data inserted successfully")
}

func seedCategories(ctx context.Context, tx pgx.Tx, logger *slog.Logger) {
	cats := []struct{ name, slug string }{
		{"Phones", "phones"},
		{"Laptops", "laptops"},
		{"Tablets", "tablets"},
	}
	for _, c := range cats {
		if _, err := tx.Exec(ctx,
			`INSERT INTO categories (name, slug) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			c.name, c.slug,
		); err != nil {
			logger.Error("seed: failed to insert category", slog.String("name", c.name), slog.String("error", err.Error()))
		}
	}
}

type productSeed struct {
	catID                                    int
	brand, model, condition                  string
	price                                    float64
	description, specs                       string
	stock                                    int
}

func seedProducts(ctx context.Context, tx pgx.Tx, sellerID int, logger *slog.Logger) {
	products := []productSeed{
		// Phones (category 1)
		{1, "Apple", "iPhone 15 Pro", "like_new", 899.99, "Titanium design with A17 Pro chip. Barely used, no scratches.", `{"storage":"256GB","ram":"8GB","color":"Natural Titanium","display":"6.1 inch","battery":"3274mAh"}`, 5},
		{1, "Apple", "iPhone 14", "excellent", 549.99, "Great condition iPhone 14 with original accessories.", `{"storage":"128GB","ram":"6GB","color":"Midnight","display":"6.1 inch","battery":"3279mAh"}`, 8},
		{1, "Apple", "iPhone 13 Mini", "good", 379.99, "Compact flagship. Minor signs of use on edges.", `{"storage":"128GB","ram":"4GB","color":"Starlight","display":"5.4 inch","battery":"2406mAh"}`, 3},
		{1, "Samsung", "Galaxy S24 Ultra", "like_new", 949.99, "S Pen included. Pristine condition with box.", `{"storage":"256GB","ram":"12GB","color":"Titanium Gray","display":"6.8 inch","battery":"5000mAh"}`, 4},
		{1, "Samsung", "Galaxy S23", "excellent", 479.99, "Snapdragon 8 Gen 2, excellent battery life.", `{"storage":"128GB","ram":"8GB","color":"Phantom Black","display":"6.1 inch","battery":"3900mAh"}`, 6},
		{1, "Google", "Pixel 8 Pro", "like_new", 699.99, "Best Android camera. Includes original packaging.", `{"storage":"128GB","ram":"12GB","color":"Obsidian","display":"6.7 inch","battery":"5050mAh"}`, 3},
		{1, "Google", "Pixel 7a", "good", 279.99, "Budget-friendly Pixel with great camera. Light wear.", `{"storage":"128GB","ram":"8GB","color":"Charcoal","display":"6.1 inch","battery":"4385mAh"}`, 10},

		// Laptops (category 2)
		{2, "Apple", "MacBook Pro 14\" M3", "like_new", 1599.99, "M3 chip, stunning Liquid Retina XDR display.", `{"storage":"512GB SSD","ram":"18GB","cpu":"Apple M3","display":"14.2 inch","battery":"70Wh"}`, 3},
		{2, "Apple", "MacBook Air 13\" M2", "excellent", 849.99, "Fanless design, perfect for daily work.", `{"storage":"256GB SSD","ram":"8GB","cpu":"Apple M2","display":"13.6 inch","battery":"52.6Wh"}`, 5},
		{2, "Lenovo", "ThinkPad X1 Carbon Gen 11", "like_new", 1199.99, "Business ultrabook, incredible keyboard.", `{"storage":"512GB SSD","ram":"16GB","cpu":"Intel i7-1365U","display":"14 inch","battery":"57Wh"}`, 4},
		{2, "Lenovo", "ThinkPad T14s", "good", 699.99, "Reliable work laptop with minor cosmetic wear.", `{"storage":"256GB SSD","ram":"16GB","cpu":"AMD Ryzen 7 7840U","display":"14 inch","battery":"52.5Wh"}`, 6},
		{2, "Dell", "XPS 15", "excellent", 1099.99, "Stunning 3.5K OLED display, powerful performance.", `{"storage":"512GB SSD","ram":"16GB","cpu":"Intel i7-13700H","display":"15.6 inch OLED","battery":"86Wh"}`, 2},
		{2, "Dell", "Latitude 5540", "fair", 549.99, "Functional business laptop, visible wear on palmrest.", `{"storage":"256GB SSD","ram":"16GB","cpu":"Intel i5-1345U","display":"15.6 inch","battery":"54Wh"}`, 7},
		{2, "HP", "EliteBook 840 G10", "excellent", 879.99, "Premium business laptop with excellent build quality.", `{"storage":"512GB SSD","ram":"32GB","cpu":"Intel i7-1365U","display":"14 inch","battery":"51Wh"}`, 3},

		// Tablets (category 3)
		{3, "Apple", "iPad Pro 12.9\" M2", "like_new", 899.99, "Pro tablet with M2 chip. Perfect for creative work.", `{"storage":"256GB","ram":"8GB","chip":"Apple M2","display":"12.9 inch","battery":"40.88Wh"}`, 4},
		{3, "Apple", "iPad Air 5th Gen", "excellent", 449.99, "M1 chip in a lightweight package.", `{"storage":"64GB","ram":"8GB","chip":"Apple M1","display":"10.9 inch","battery":"28.6Wh"}`, 6},
		{3, "Apple", "iPad 10th Gen", "good", 299.99, "Redesigned iPad with USB-C. Minor screen scratch.", `{"storage":"64GB","ram":"4GB","chip":"A14 Bionic","display":"10.9 inch","battery":"28.6Wh"}`, 8},
		{3, "Samsung", "Galaxy Tab S9", "like_new", 599.99, "AMOLED display with S Pen. Like new with box.", `{"storage":"128GB","ram":"8GB","chip":"Snapdragon 8 Gen 2","display":"11 inch","battery":"8400mAh"}`, 3},
		{3, "Samsung", "Galaxy Tab A9+", "excellent", 219.99, "Affordable tablet, great for media consumption.", `{"storage":"64GB","ram":"4GB","chip":"Snapdragon 695","display":"11 inch","battery":"7040mAh"}`, 12},
		{3, "Lenovo", "Tab P12 Pro", "good", 349.99, "AMOLED display, quad speakers. Light scratches on back.", `{"storage":"256GB","ram":"8GB","chip":"Snapdragon 870","display":"12.6 inch","battery":"10200mAh"}`, 5},
	}

	for _, p := range products {
		if _, err := tx.Exec(ctx,
			`INSERT INTO products (seller_id, category_id, brand, model, condition, price, description, image_url, specs, stock)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, '', $8::jsonb, $9)`,
			sellerID, p.catID, p.brand, p.model, p.condition, p.price, p.description, p.specs, p.stock,
		); err != nil {
			logger.Error("seed: failed to insert product", slog.String("model", p.model), slog.String("error", err.Error()))
		}
	}
}
