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

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"notification-service/internal/handler"
	"notification-service/internal/middleware"
	"notification-service/internal/repository"
	"notification-service/internal/service"
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

	repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(repo)
	h := handler.NewNotificationHandler(svc)

	authService := service.NewAuthService(jwtSecret)
	authMiddleware := middleware.JWTAuthMiddleware(authService)

	// routing
	namedTypedNotificationHandler := handler.AppHandler(h.HandleNotifications)

	// Combine logging and auth middlewares
	loggedHandler := middleware.LoggingMiddleware(namedTypedNotificationHandler)
	protectedHandler := authMiddleware(loggedHandler)

	http.Handle("/notifications", protectedHandler)

	// TODO(лекция про Kafka): здесь нужно поднять Kafka consumer, который
	// слушает топик "order.created" (его публикует order-service после
	// успешного создания заказа) и на каждое сообщение дергает
	// svc.CreateNotification(ctx, event.UserID, event.OrderID).
	// Сам consumer в проект намеренно не добавлен.

	server := http.Server{
		Addr: ":8083",
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
