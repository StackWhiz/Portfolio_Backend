package database

import (
	"arbak-portfolio-backend/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize sets up the database connection and runs migrations
func Initialize(databaseURL string) (*gorm.DB, error) {
	// Configure GORM logger
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(databaseURL), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed initial data if needed
	if err := seedInitialData(db); err != nil {
		log.Printf("Warning: failed to seed initial data: %v", err)
	}

	return db, nil
}

// InitializeRedis sets up Redis connection
func InitializeRedis(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Warning: failed to parse Redis URL, using default config: %v", err)
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: failed to connect to Redis: %v", err)
	}

	return client
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Profile{},
		&models.Experience{},
		&models.Skill{},
		&models.Project{},
		&models.Contact{},
		&models.User{},
	)
}

// seedInitialData seeds the database with initial data
func seedInitialData(db *gorm.DB) error {
	// Check if profile already exists
	var count int64
	db.Model(&models.Profile{}).Count(&count)
	if count > 0 {
		return nil // Data already exists
	}

	// Create initial profile
	profile := &models.Profile{
		Name:     "Your name",
		Title:    "title",
		Location: "location",
		Email:    "email@gmail.com",
		Phone:    "+123456789",
		Telegram: "@telegram",
		GitHub:   "github.com/StackWhiz",
		Summary:  `summary.`,
	}

	if err := db.Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create initial profile: %w", err)
	}

	// Create initial experiences
	experiences := []models.Experience{
		{
			Company:     "Company1",
			Position:    "Position",
			Location:    "Remote",
			StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Current:     true,
			Description: "Description",
			Achievements: []string{
				"Architected and led backend services in Rust and Go, scaling APIs and microservices to handle millions of daily requests",
				"Implemented PoS consensus logic and validator services in Rust, enhancing block finality and network reliability",
				"Built Kafka + Postgres + ClickHouse pipelines processing 50k+ blockchain events per second",
				"Developed and audited Solidity & Anchor smart contracts for staking, governance, token bridging, and liquidity pools",
				"Designed DDoS protection strategies (rate-limiting, WAF, caching, load balancing) securing validator RPCs and public APIs",
				"Containerized workloads with Docker and deployed to Kubernetes (GKE) with Helm, Prometheus/Grafana, and ELK logging",
				"Established CI/CD pipelines (GitHub Actions + GitLab CI) automating builds, tests, and deployments",
				"Led and mentored 6 engineers, introducing best practices in distributed systems, DevOps, and blockchain protocol design",
			},
			Technologies: []string{"Rust", "Go", "Kafka", "PostgreSQL", "ClickHouse", "Solidity", "Anchor", "Docker", "Kubernetes", "Helm", "Prometheus", "Grafana"},
		},
		{
			Company:     "Company2",
			Position:    "Position",
			Location:    "Remote",
			StartDate:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     &[]time.Time{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}[0],
			Current:     false,
			Description: "Developed high-performance trading systems and secure wallet infrastructure",
			Achievements: []string{
				"Developed and optimized a Go-based matching engine sustaining 10k+ TPS with <50ms latency",
				"Designed and deployed trading APIs (REST, WebSocket, gRPC) serving 50k+ concurrent users",
				"Built secure wallet microservices in Rust with multi-sig and HSM integrations",
				"Architected DDoS-resistant API gateways with throttling, reverse proxies, and auto-scaling clusters",
				"Optimized PostgreSQL sharding and Redis caching, boosting performance by 35%",
				"Automated deployments with CI/CD pipelines (Docker + GitLab CI), reducing release times by 60%",
				"Delivered 99.99% uptime SLA across multi-region Kubernetes clusters (AWS & GCP)",
				"Contributed to MEV-resistant order execution logic, mitigating front-running attacks",
			},
			Technologies: []string{"Go", "Rust", "PostgreSQL", "Redis", "Docker", "Kubernetes", "AWS", "GCP", "gRPC", "WebSocket"},
		},
		{
			Company:     "Company3",
			Position:    "Position",
			Location:    "Remote",
			StartDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     &[]time.Time{time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)}[0],
			Current:     false,
			Description: "Built blockchain analytics and transaction indexing systems",
			Achievements: []string{
				"Built Rust & Go-based microservices for transaction indexing and real-time blockchain analytics",
				"Implemented fraud/anomaly detection modules with Kafka + ClickHouse, improving detection accuracy by 20%",
				"Developed GraphQL + REST APIs serving blockchain insights to enterprise clients",
				"Designed streaming architectures with Kafka, ClickHouse, and Redis, enabling <1s latency dashboards",
				"Enhanced node protocols for mempool data capture and transaction propagation, improving throughput by 30%",
				"Containerized applications with Docker and set up automated pipelines for staging/production",
			},
			Technologies: []string{"Rust", "Go", "Kafka", "ClickHouse", "Redis", "GraphQL", "Docker"},
		},
		{
			Company:     "Company4",
			Position:    "Position",
			Location:    "Remote",
			StartDate:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     &[]time.Time{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}[0],
			Current:     false,
			Description: "Developed financial transaction processing systems",
			Achievements: []string{
				"Developed Go microservices handling 100k+ daily financial transactions",
				"Integrated ISO8583 and SWIFT protocols, ensuring compliance with global banking standards",
				"Built fraud detection engines using Redis + Postgres triggers, reducing fraudulent cases by 25%",
				"Designed secure API gateways with JWT auth, rate-limiting, and RBAC",
				"Implemented DDoS protection layers with load balancing and request filtering",
				"Automated compliance reporting workflows, cutting audit effort by 40%",
			},
			Technologies: []string{"Go", "PostgreSQL", "Redis", "JWT", "ISO8583", "SWIFT"},
		},
	}

	for _, exp := range experiences {
		if err := db.Create(&exp).Error; err != nil {
			return fmt.Errorf("failed to create experience: %w", err)
		}
	}

	// Create initial skills
	skills := []models.Skill{
		// Languages
		{Name: "Rust", Category: "Languages", Level: 9, Description: "Systems programming, blockchain development", Icon: "🦀"},
		{Name: "Go", Category: "Languages", Level: 9, Description: "Backend services, microservices", Icon: "🐹"},
		{Name: "JavaScript/TypeScript", Category: "Languages", Level: 8, Description: "Full-stack development", Icon: "🟨"},
		{Name: "Python", Category: "Languages", Level: 7, Description: "Data processing, automation", Icon: "🐍"},
		{Name: "Solidity", Category: "Languages", Level: 8, Description: "Smart contract development", Icon: "⛓️"},

		// Frameworks
		{Name: "Actix", Category: "Frameworks", Level: 8, Description: "Rust web framework", Icon: "⚡"},
		{Name: "Axum", Category: "Frameworks", Level: 7, Description: "Rust async web framework", Icon: "🪶"},
		{Name: "Echo", Category: "Frameworks", Level: 8, Description: "Go web framework", Icon: "🌊"},
		{Name: "Gin", Category: "Frameworks", Level: 8, Description: "Go HTTP web framework", Icon: "🍸"},
		{Name: "Express.js", Category: "Frameworks", Level: 7, Description: "Node.js web framework", Icon: "🚀"},
		{Name: "NestJS", Category: "Frameworks", Level: 7, Description: "Node.js enterprise framework", Icon: "🏗️"},

		// Blockchain
		{Name: "Consensus Algorithms", Category: "Blockchain", Level: 9, Description: "PoS, BFT consensus implementation", Icon: "🔗"},
		{Name: "Validator Nodes", Category: "Blockchain", Level: 9, Description: "Blockchain validator infrastructure", Icon: "⚖️"},
		{Name: "MEV & DeFi", Category: "Blockchain", Level: 8, Description: "MEV infrastructure, DeFi protocols", Icon: "💰"},
		{Name: "P2P Networking", Category: "Blockchain", Level: 8, Description: "Distributed network protocols", Icon: "🌐"},

		// DevOps
		{Name: "Docker", Category: "DevOps", Level: 9, Description: "Containerization", Icon: "🐳"},
		{Name: "Kubernetes", Category: "DevOps", Level: 8, Description: "Container orchestration", Icon: "☸️"},
		{Name: "Helm", Category: "DevOps", Level: 7, Description: "Kubernetes package manager", Icon: "⛵"},
		{Name: "AWS", Category: "DevOps", Level: 8, Description: "Cloud infrastructure", Icon: "☁️"},
		{Name: "Azure", Category: "DevOps", Level: 7, Description: "Microsoft cloud platform", Icon: "🔷"},

		// Databases
		{Name: "PostgreSQL", Category: "Databases", Level: 9, Description: "Relational database", Icon: "🐘"},
		{Name: "Redis", Category: "Databases", Level: 8, Description: "In-memory data store", Icon: "🔴"},
		{Name: "ClickHouse", Category: "Databases", Level: 7, Description: "Analytical database", Icon: "📊"},
		{Name: "MongoDB", Category: "Databases", Level: 6, Description: "NoSQL document database", Icon: "🍃"},
		{Name: "Cassandra", Category: "Databases", Level: 6, Description: "Distributed NoSQL database", Icon: "🗃️"},
	}

	for _, skill := range skills {
		if err := db.Create(&skill).Error; err != nil {
			return fmt.Errorf("failed to create skill: %w", err)
		}
	}

	// Create initial projects
	projects := []models.Project{
		{
			Name:            "High-Performance Trading Engine",
			Description:     "Go-based matching engine sustaining 10k+ TPS with <50ms latency",
			LongDescription: "Built a high-frequency trading engine using Go with custom data structures and memory optimization techniques. Implemented order matching algorithms, real-time market data distribution, and risk management systems.",
			Technologies:    []string{"Go", "Redis", "PostgreSQL", "WebSocket", "gRPC"},
			Category:        "Backend",
			Featured:        true,
			Status:          "completed",
		},
		{
			Name:            "Blockchain Validator Infrastructure",
			Description:     "Rust-based validator services with PoS consensus implementation",
			LongDescription: "Developed and deployed blockchain validator infrastructure using Rust. Implemented custom consensus algorithms, P2P networking protocols, and monitoring systems for high availability.",
			Technologies:    []string{"Rust", "Docker", "Kubernetes", "Prometheus", "Grafana"},
			Category:        "Blockchain",
			Featured:        true,
			Status:          "completed",
		},
		{
			Name:            "Real-time Analytics Pipeline",
			Description:     "Kafka + ClickHouse pipeline processing 50k+ blockchain events per second",
			LongDescription: "Architected a real-time data processing pipeline for blockchain analytics. Built streaming data ingestion, real-time aggregation, and dashboard systems for enterprise clients.",
			Technologies:    []string{"Kafka", "ClickHouse", "Rust", "Go", "Redis"},
			Category:        "Backend",
			Featured:        true,
			Status:          "completed",
		},
		{
			Name:            "Smart Contract Suite",
			Description:     "Solidity & Anchor smart contracts for DeFi protocols",
			LongDescription: "Developed comprehensive smart contract suite including staking mechanisms, governance systems, token bridging protocols, and liquidity pools with security audits.",
			Technologies:    []string{"Solidity", "Anchor", "Rust", "TypeScript"},
			Category:        "Blockchain",
			Featured:        true,
			Status:          "completed",
		},
	}

	for _, project := range projects {
		if err := db.Create(&project).Error; err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
	}

	return nil
}
