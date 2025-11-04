# MarketFlow

MarketFlow is a real-time market data processing system built in Go using Hexagonal Architecture. The application collects data from cryptocurrency exchange simulators or generates test data, aggregates prices, stores them in PostgreSQL, and caches them in Redis. A built-in REST API provides convenient access to aggregated market information.

## 🚀 Live Demo

View Live Application
(replace with actual link if available)

## 🛠️ Technologies Used

Backend: Go (1.21+)

Database: PostgreSQL

Cache: Redis

Deployment: Docker, Docker Compose

## ✨ Features

Real-time aggregation of market prices

Live/Test modes for flexible data sources

Worker pool for concurrent feed processing (5 workers per exchange)

Fan-In / Fan-Out architecture for data streams

Batch insertion into PostgreSQL for efficiency

Automatic fallback to DB if Redis is unavailable

REST API to fetch latest, highest, lowest, and average prices

System health endpoint and structured logging

## 📦 Installation

Clone the repository:

git clone hgit@github.com:AikaBe/LeaderBoard-System.git
cd marketflow


Load exchange simulator images (if needed):

docker-compose run --rm load_images


Build the project:

docker-compose build


Configure config.yaml:

postgres:
host: localhost
port: 5432
user: marketflow
password: secret
dbname: marketflow_db

redis:
host: localhost
port: 6379
password: ""

exchanges:
- name: exchange1
  host: 127.0.0.1
  port: 40101
- name: exchange2
  host: 127.0.0.1
  port: 40102
- name: exchange3
  host: 127.0.0.1
  port: 40103

## 🎯 Usage
Start the application with Docker Compose:
docker-compose up

Running exchange simulators (Live Mode):
docker load -i exchange1_amd64.tar
docker run -p 40101:40101 -d exchange1_amd64

docker load -i exchange2_amd64.tar
docker run -p 40102:40102 -d exchange2_amd64

docker load -i exchange3_amd64.tar
docker run -p 40103:40103 -d exchange3_amd64

API Examples

Fetch the latest price for BTCUSDT:

curl http://localhost:8080/prices/latest/BTCUSDT


Switch to test mode:

curl -X POST http://localhost:8080/mode/test


Check system health:

curl http://localhost:8080/health

## 🏗️ Architecture

Hexagonal Architecture (Ports & Adapters):

Domain Layer: business logic (price aggregation, data models)

Application Layer: use-case processing, worker pool management, data flow handling

Adapters:

HTTP Adapter (REST API)

Storage Adapter (PostgreSQL)

Cache Adapter (Redis)

Exchange Adapter (Live/Test sources)

The system supports fan-in/fan-out data processing, batch inserts, Redis caching with fallback, and automatic reconnection to data sources in case of failure.

## 🔮 Future Improvements

Add WebSocket endpoints for live price streaming

Support additional trading pairs and exchanges

Add authentication and API key management

Implement historical data analytics and visualization