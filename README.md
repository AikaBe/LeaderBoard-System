# LeaderBoard
Leaderboard is a simplified anonymous imageboard inspired by early internet forums. Users can create posts, comment, and share images without the need for user registration. This project is built using Go and PostgreSQL with a focus on Hexagonal Architecture, session management, and integration with external services like The Rick and Morty API for user avatars.

## 🛠️ Technologies Used

Frontend: Provided HTML templates (customizable, no frontend development required)

Backend: Go (Golang), PostgreSQL, S3-compatible storage (MinIO or similar)

Architecture: Hexagonal Architecture (Ports and Adapters pattern)

Session Management: HTTP cookies (no user registration)

External APIs: Rick and Morty API for user avatars

Deployment: Docker (for local development)

## ✨ Features

Anonymous posts and comments: Users can create posts with text, images, and comments. No registration required.

Image uploads: Attach images to posts and comments, stored securely on S3-compatible storage.

User avatars: Unique avatars assigned to users using the Rick and Morty API.

Session-based user identification: Each session is tracked via cookies, ensuring a persistent user experience.

Post expiration: Posts without comments are deleted after 10 minutes; posts with comments are deleted after 15 minutes of inactivity.

Responsive design: Mobile and desktop-friendly design with provided templates.

RESTful API: Backend server exposes a REST API that interacts with the frontend.

## 📦 Installation
Requirements

Go (v1.16+)

PostgreSQL database

S3-compatible storage (MinIO or alternative)

Setup Instructions

Clone the repository:

git clone https://github.com/AikaBe/LeaderBoard-System.git
cd 1337b04rd


Install Go dependencies:

go mod tidy


Set up the PostgreSQL database:
Create a PostgreSQL database and configure the connection details in the .env file.

Example .env file:

DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=1337b04rd


Start S3-compatible storage (e.g., MinIO):

Run MinIO locally or use an S3 provider for image storage.

Set the appropriate environment variables for storage credentials in the .env file.

Run the server:

go run ./cmd/1337b04rd


Access the application:
Navigate to http://localhost:8080 in your browser to interact with the imageboard.

## 🎯 Usage
Starting the Application
./1337b04rd --port 8080


This will start the backend server. You can access the imageboard via http://localhost:8080.

Post and Comment Creation

Create a new post: Visit the main page and click on "Create New Post". You can add text and optionally attach an image.

Add comments: To comment on a post, navigate to the post's page and type your comment. You can reply to specific comments by clicking on their ID.

View posts: On the main page (catalog.html), you'll see active threads, and archived threads can be accessed via the "Archive" button.

Session Management

The system uses cookies to track user sessions. Upon the first visit, each user is assigned a unique avatar and name from the Rick and Morty API.

## 🏗️ Architecture
Hexagonal Architecture

This project follows Hexagonal Architecture (also known as Ports and Adapters pattern), which separates core business logic from external systems (e.g., database, image storage, web server). This ensures that the core functionality remains independent of how data is stored or served.

Domain Layer:

Contains core business logic for creating posts, comments, and managing sessions.

Defines interfaces (ports) for data storage and external services.

Infrastructure Layer:

Concrete implementations of the domain interfaces (e.g., PostgreSQL for data storage, MinIO for image storage, and external APIs for avatars).

User Interface Layer:

Handles incoming HTTP requests, session management, and routing.

Serves the provided HTML templates to interact with the user.

Session and User Identification

Users are tracked via cookies. On the first visit, they are assigned a unique avatar and name fetched from the Rick and Morty API.

Avatars are unique for each user per session. If all avatars are used, they may be reused.

## 🔮 Future Improvements

Admin features: Implement an admin interface to manage posts and comments manually.

Search functionality: Add the ability to search for posts or comments by keywords.

User registration: Optional feature to enable user registration for persistent profiles.

Advanced error handling: Improve error messages and response handling for users.
