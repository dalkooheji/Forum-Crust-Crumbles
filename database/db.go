package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() error {
	dbPath := filepath.Join("database", "forum.db")

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Database does not exist. Creating...")
		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("error creating database file: %v", err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	DB = db

	if err := createTables(); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	if isEmpty, err := isDatabaseEmpty(); err != nil {
		return fmt.Errorf("error checking if database is empty: %v", err)
	} else if isEmpty {
		if err := seedData(); err != nil {
			return fmt.Errorf("error seeding data: %v", err)
		}
	}

	fmt.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	sqlFile, err := os.ReadFile("database/db_structure.sql")
	if err != nil {
		return fmt.Errorf("error reading SQL file: %v", err)
	}

	// Split the SQL file content into individual statements
	statements := strings.Split(string(sqlFile), ";")

	for _, stmt := range statements {
		// Trim whitespace and skip empty statements
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// Add IF NOT EXISTS to CREATE TABLE statements
		if strings.HasPrefix(strings.ToUpper(stmt), "CREATE TABLE") {
			stmt = strings.Replace(stmt, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS", 1)
		}

		_, err = DB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing SQL: %v", err)
		}
	}

	return nil
}
func isDatabaseEmpty() (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM Users").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
func seedData() error {
	rand.Seed(time.Now().UnixNano())

	categories := []string{"Simple", "Advanced", "Appetizers", "Main", "Desserts", "Drinks"}
	users := []string{"alice", "bob", "charlie", "david", "eve", "frank", "grace", "henry", "isabel", "jack"}
	foodPosts := []struct {
		Title   string
		Content string
	}{
		{"Creamy Garlic Parmesan Pasta", "Ingredients: pasta, heavy cream, garlic, parmesan cheese. Cook pasta, sauté garlic, add cream and cheese, toss with pasta."},
		{"Perfect Grilled Steak", "Season ribeye with salt and pepper. Grill 4-5 minutes per side for medium-rare. Let rest before serving."},
		{"Quick Avocado Toast", "Toast bread, mash avocado with lemon juice and salt. Spread on toast and top with a poached egg."},
		{"Decadent Chocolate Lava Cake", "Mix chocolate, butter, sugar, eggs, and flour. Bake in ramekins for 12 minutes at 400°F for a gooey center."},
		{"Cold Brew Coffee", "Coarsely grind coffee beans. Steep in cold water for 12 hours. Strain and serve over ice."},
		{"Green Detox Smoothie", "Blend spinach, banana, apple, ginger, and coconut water for a refreshing and healthy drink."},
		{"Homemade Margherita Pizza", "Make dough, top with San Marzano tomatoes, fresh mozzarella, and basil. Bake in a hot oven until crispy."},
		{"Aromatic Indian Curry", "Sauté onions, add spices, tomatoes, and coconut milk. Simmer with your choice of protein or vegetables."},
		{"Sushi Rolling Basics", "Prepare sushi rice, lay nori sheet, add rice and fillings. Roll tightly and slice into pieces."},
		{"Comforting Chicken Noodle Soup", "Simmer chicken, carrots, celery, and egg noodles in broth. Season with thyme and parsley."},
		{"Crispy Baked Sweet Potato Fries", "Cut sweet potatoes into wedges, toss with oil and spices. Bake until crispy, serve with aioli."},
		{"Refreshing Greek Salad", "Mix cucumbers, tomatoes, red onion, olives, and feta. Dress with olive oil and oregano."},
		{"Homemade Guacamole", "Mash avocados with lime juice, diced onion, tomato, cilantro, and jalapeño. Season to taste."},
		{"Fluffy Pancakes from Scratch", "Mix flour, baking powder, milk, eggs, and melted butter. Cook on a griddle until golden brown."},
		{"Spicy Vegetarian Chili", "Simmer beans, tomatoes, bell peppers, and onions with chili spices. Top with cheese and sour cream."},
		{"Classic French Onion Soup", "Caramelize onions slowly, add beef broth and wine. Top with bread and gruyere, broil until cheese melts."},
		{"Homemade Pesto Sauce", "Blend fresh basil, pine nuts, garlic, parmesan, and olive oil. Perfect for pasta or as a spread."},
		{"Crispy Fried Chicken", "Marinate chicken in buttermilk, coat in seasoned flour. Fry until golden and crispy."},
		{"Vegetable Stir Fry", "Quick-fry assorted vegetables in a wok with soy sauce and ginger. Serve over rice."},
		{"Creamy Mushroom Risotto", "Slowly add broth to arborio rice, stir in sautéed mushrooms and parmesan for a luxurious dish."},
	}

	// Create users
	for _, username := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}

		_, err = DB.Exec("INSERT INTO Users (Username, Email, PasswordHash) VALUES (?, ?, ?)",
			username, username+"@example.com", string(hashedPassword))
		if err != nil {
			return fmt.Errorf("error inserting user: %v", err)
		}
	}

	// Create categories
	for _, category := range categories {
		_, err := DB.Exec("INSERT INTO Categories (CategoryName) VALUES (?)", category)
		if err != nil {
			return fmt.Errorf("error inserting category: %v", err)
		}
	}

	// Create posts, likes, dislikes, and comments
	for i, user := range users {
		var userID int
		err := DB.QueryRow("SELECT UserID FROM Users WHERE Username = ?", user).Scan(&userID)
		if err != nil {
			return fmt.Errorf("error getting user ID: %v", err)
		}

		for j := 0; j < 2; j++ {
			postIndex := (i*2 + j) % len(foodPosts)
			post := foodPosts[postIndex]

			var postID int64
			result, err := DB.Exec("INSERT INTO Posts (UserID, Title, Content) VALUES (?, ?, ?)",
				userID, post.Title, post.Content)
			if err != nil {
				return fmt.Errorf("error inserting post: %v", err)
			}
			postID, _ = result.LastInsertId()

			// Assign random category to post
			categoryID := rand.Intn(len(categories)) + 1
			_, err = DB.Exec("INSERT INTO PostCategories (PostID, CategoryID) VALUES (?, ?)", postID, categoryID)
			if err != nil {
				return fmt.Errorf("error assigning category to post: %v", err)
			}

			// Create likes and dislikes
			for k := 0; k < 2; k++ {
				likerID := rand.Intn(len(users)) + 1
				_, err = DB.Exec("INSERT INTO Likes (UserID, PostID, IsLike) VALUES (?, ?, ?)", likerID, postID, true)
				if err != nil {
					return fmt.Errorf("error inserting like: %v", err)
				}

				dislikerID := rand.Intn(len(users)) + 1
				_, err = DB.Exec("INSERT INTO Dislikes (UserID, PostID, IsDislike) VALUES (?, ?, ?)", dislikerID, postID, true)
				if err != nil {
					return fmt.Errorf("error inserting dislike: %v", err)
				}
			}

			// Create comments
			for k := 0; k < 2; k++ {
				commenterID := rand.Intn(len(users)) + 1
				commentContent := fmt.Sprintf("This is comment %d on the post about %s", k+1, post.Title)
				_, err = DB.Exec("INSERT INTO Comments (PostID, UserID, Content) VALUES (?, ?, ?)",
					postID, commenterID, commentContent)
				if err != nil {
					return fmt.Errorf("error inserting comment: %v", err)
				}
			}
		}
	}

	fmt.Println("Fake data seeded successfully")
	return nil
}
