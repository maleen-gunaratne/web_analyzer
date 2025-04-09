## Go Web Analyzer

# Project Overview

Go Web Analyzer is a web service that analyzes HTML pages, extracts useful metadata, and provides insights such as 
HTML version, title, headings, internal and external links, broken links, and login form detection. 
It is built using Go and the Gin framework, with Prometheus integration for monitoring.

Before running the project, ensure you have the following installed:

Go (v1.18+)

Gin (Web framework)

Prometheus (Metrics monitoring)


-Technology Stack

Backend:

 Language: Go (Golang)

 Framework: Gin

Frontend :

 Language : JavaScript

 Framework : React

Features

- HTML version detection
- Page title extraction
- Heading structure analysis
- Internal and external link identification
- Broken link detection
- Login form detection
- Prometheus metrics integration

------------------------------------------------------------------------------

# Installation and Setup (Backend - Go framework)

Clone the Repository

git clone https://github.com/maleen-gunaratne/web_analyzer.git

cd web-analyzer

Ensure you have the required Go dependencies installed:

 - Install Go modules

go mod tidy

Running the Application

- Run the server

go run cmd/main.go

The server will start on port 8080. You can access the API at
http://localhost:8080/url-analyze.

---------------------------------------------------------------------------------------------

# Build and run the Docker container

docker build -t web-analyzer .

docker run -p 8080:8080 web-analyzer
