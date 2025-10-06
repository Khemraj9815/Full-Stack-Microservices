// services/products-service/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "practical-three/proto/gen"
	consulapi "github.com/hashicorp/consul/api"
)

const serviceName = "products-service"
const servicePort = 50052

type Product struct {
	gorm.Model
	Name  string
	Price float64
}

type server struct {
	pb.UnimplementedProductServiceServer
	db *gorm.DB
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product := Product{Name: req.Name, Price: req.Price}
	if result := s.db.Create(&product); result.Error != nil {
		return nil, result.Error
	}
	return &pb.ProductResponse{Product: &pb.Product{Id: fmt.Sprint(product.ID), Name: product.Name, Price: product.Price}}, nil
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	var product Product
	if result := s.db.First(&product, req.Id); result.Error != nil {
		return nil, result.Error
	}
	return &pb.ProductResponse{Product: &pb.Product{Id: fmt.Sprint(product.ID), Name: product.Name, Price: product.Price}}, nil
}

func main() {
	// Wait for the database to be ready
	time.Sleep(5 * time.Second)

	// Connect to the database
	dsn := "host=products-db user=user password=password dbname=products_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Product{})

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{db: db})

	// Register with Consul after gRPC server is listening
	go func() {
		time.Sleep(2 * time.Second)
		if err := registerServiceWithConsul(); err != nil {
			log.Fatalf("Failed to register with Consul: %v", err)
		}
	}()

	log.Printf("%s gRPC server listening at %v", serviceName, lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func registerServiceWithConsul() error {
	config := consulapi.DefaultConfig()
	if addr := os.Getenv("CONSUL_HTTP_ADDR"); addr != "" {
		config.Address = addr
	}

	consul, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%d", serviceName, servicePort),
		Name:    serviceName,
		Port:    servicePort,
		Address: serviceName, // Docker service name
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", serviceName, servicePort),
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	return consul.Agent().ServiceRegister(registration)
}
