package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"order-service/internal/client"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"order-service/internal/handler"
	"order-service/internal/middleware"
	"order-service/internal/repository"
	"order-service/internal/service"
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

	productClient := client.NewProductClient("http://localhost:8080")

	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo, productClient)
	h := handler.NewOrderHandler(svc)

	authService := service.NewAuthService(jwtSecret)
	authMiddleware := middleware.JWTAuthMiddleware(authService)

	// routing
	namedTypedOrderHandler := handler.AppHandler(h.HandleOrders)

	// Combine logging and auth middlewares
	loggedHandler := middleware.LoggingMiddleware(namedTypedOrderHandler)
	protectedHandler := authMiddleware(loggedHandler)

	http.Handle("/orders", protectedHandler)

	server := http.Server{
		Addr: ":8082",
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
