# Gemini Code Assist Project Overview: Wiki-Voyage Demo

This document provides an overview of the Wiki-Voyage demo application, its structure, coding standards, and CI/CD pipeline, specifically tailored for analysis and interaction with Gemini Code Assist.

## About This Project

This project is a web application written in Go. It serves as a demonstration of how to build a simple microservice that can be deployed to Google Cloud Run. The application provides two main endpoints:

*   `/`: Lists points of interest.
*   `/recommendations`: Provides recommendations.

The application is designed to be containerized and deployed using a CI/CD pipeline.

## Project Structure

The project is structured as follows:

```
├───.gcloudignore
├───.gitignore
├───bitbucket-pipelines.yml
├───cloudbuild.yaml
├───Dockerfile
├───go.mod
├───go.sum
├───main.go
├───models.go
├───poilist.go
├───README.md
├───recommendations.go
├───.git/...
├───templates/
│   └───index.html
└───tmp/...
```

*   **`main.go`**: The main entry point for the application. It starts the web server and registers the HTTP handlers.
*   **`poilist.go`**: Contains the logic for the `/` endpoint, which lists points of interest.
*   **`recommendations.go`**: Contains the logic for the `/recommendations` endpoint.
*   **`models.go`**: Defines the data structures used in the application.
*   **`templates/index.html`**: The HTML template for the main page.
*   **`go.mod`**, **`go.sum`**: Go module files that manage the project's dependencies.
*   **`Dockerfile`**: Used to build the Docker container for the application.
*   **`cloudbuild.yaml`**: The configuration file for Google Cloud Build, which builds the Docker image and deploys it to Cloud Run.
*   **`bitbucket-pipelines.yml`**: The configuration file for Bitbucket Pipelines, which is set up to perform automated code reviews on pull requests using the Gemini CLI.

## Go Coding Standards

This project follows standard Go coding practices:

*   **Formatting**: All Go code is formatted using `gofmt`.
*   **Linting**: `go vet` is used to identify potential issues in the code.
*   **Dependencies**: Dependencies are managed using Go modules.

## CI/CD

The project has a CI/CD pipeline that uses Bitbucket Pipelines and Google Cloud Build.

### Bitbucket Pipelines

The `bitbucket-pipelines.yml` file defines a pipeline that runs on every pull request. The pipeline uses the Gemini CLI to perform an automated code review of the changes in the pull request. The results of the review are then posted as a comment on the pull request in Bitbucket.

### Google Cloud Build

The `cloudbuild.yaml` file defines a pipeline that:

1.  Builds a Docker image of the application using Kaniko.
2.  Pushes the image to Google Artifact Registry.
3.  Deploys the image to Google Cloud Run.

This pipeline is triggered manually or by changes to the main branch.

## How to Build and Run

To build and run the application locally, you can use the following commands:

```bash
# Build the application
go build -o wiki-voyage

# Run the application
./wiki-voyage
```

The application will be available at `http://localhost:8084`.

### Hot Reloading with Air

For a better development experience with hot reloading, you can use `air`. `air` will automatically rebuild and restart the application when it detects changes in the source code.

To install `air`:

```bash
go install github.com/cosmtrek/air@latest
```

To run the application with `air`:

```bash
air
```

This will start the application and watch for any changes in your Go files.