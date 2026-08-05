package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"user-service/internal/handler"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/service"
)

var db *sql.DB
var jwtSecret = []byte("sultanbekJwtKey")

func InitDB() {
	connection := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	errTwo := db.Ping()
	if errTwo != nil {
		log.Println("Warning: Failed to ping database: ", errTwo)
	}
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitDB()
	defer CloseDB()

	// Run Migrations
	runMigrations(db, "file://migrations")

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepository, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	// routing
	namedTypedUserHandler := handler.AppHandler(userHandler.HandleUsers)

	http.Handle("/users", middleware.LoggingMiddleware(namedTypedUserHandler))
	http.Handle("/login", http.HandlerFunc(authHandler.Login))
	http.Handle("/logout", http.HandlerFunc(authHandler.Logout))

	server := http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("server started on", server.Addr)
		// По стандартам Go нужно игнорировать ErrServerClosed, так как это ожидаемое поведение при Shutdown
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// создали канал который будет из горутины ос ---> принимать сигналы
	// По стандартам канал для сигналов должен быть буферизированным (размер 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// считываю из канала
	<-quit

	log.Println("Считывание из канала закончено")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Вызываем shutdown")
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("server has been shutdown with graceful 5s error: ", err)
	}
}

func runMigrations(db *sql.DB, sourceURL string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Migration warning: failed to create postgres driver instance (expected if db is down): %v", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		log.Printf("Migration warning: failed to initialize migrate instance: %v", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration warning: failed to run migrate up: %v", err)
		return
	}

	log.Println("Migrations applied successfully!")
}
