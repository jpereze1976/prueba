import (
	"context"
	"fmt"
	"log"
	"os"

	"solarwinds-backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Cargar .env
	_ = godotenv.Load()

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Conectar a la base de datos
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Datos del usuario
	email := os.Getenv("USER_EMAIL")
	password := os.Getenv("USER_PASSWORD")
	fullName := os.Getenv("USER_FULLNAME")
	role := os.Getenv("USER_ROLE")

	if email == "" || password == "" || fullName == "" {
		fmt.Println("Usage:")
		fmt.Println("  USER_EMAIL=admin@example.com USER_PASSWORD=admin123 USER_FULLNAME=\"Admin User\" USER_ROLE=admin go run main.go")
		fmt.Println("")
		fmt.Println("Roles disponibles: admin, user")
		fmt.Println("Si no se especifica USER_ROLE, se usa 'admin' por defecto")
		os.Exit(1)
	}

	// Si no se especifica rol, usar 'admin' por defecto
	if role == "" {
		role = "admin"
	}

	// Validar rol
	if role != "admin" && role != "user" {
		fmt.Println("Error: USER_ROLE debe ser 'admin' o 'user'")
		os.Exit(1)
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Insertar usuario
	query := `
		INSERT INTO users (email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id
	`

	var userID string
	err = db.QueryRow(ctx, query, email, string(hashedPassword), fullName, role).Scan(&userID)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}

	fmt.Printf("✅ User created successfully!\n")
	fmt.Printf("   ID: %s\n", userID)
	fmt.Printf("   Email: %s\n", email)
	fmt.Printf("   Role: %s\n", role)
	fmt.Printf("\n🔑 You can now login with:\n")
	fmt.Printf("   Email: %s\n", email)
	fmt.Printf("   Password: %s\n", password)
}
