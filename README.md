# Wiki Voyage Demo Application

This application demonstrates the capabilities of Gemini Code Assist in a real-world coding scenario. It's a simple web application that displays points of interest from a BigQuery dataset, showcasing various features like data loading, code generation, testing, refactoring, and security analysis with Gemini Code Assist's assistance.

## Functionality

The application fetches points of interest for a specified city from a BigQuery dataset (`wiki_voyage.points_of_interest`). It then presents this data in a user-friendly format on a web page.  Additional features include:

* **Data Loading from BigQuery:**  The application dynamically loads data from BigQuery based on user input.
* **Emoji Representation of Activities:**  Activity types are represented by relevant emojis.
* **Recommendations:** (Future enhancement)  The application will eventually generate recommendations based on the points of interest.

## Getting Started

This demo requires several setup steps before running. Please refer to the detailed instructions in the main `README.md` file in this repository.  Key steps include:

* **Setting up your Google Cloud Project:**  This involves enabling necessary Google Cloud services, creating a BigQuery dataset, and setting environment variables.
* **Cloning the Repository:** Clone the repository containing this application and its dependencies.
* **Installing Dependencies:** Make sure you have `go` installed and run `go install github.com/air-verse/air@latest`
* **Running the application:** Navigate to the application's directory and run `air` in your terminal. This will start a local development server with hot reload capabilities using air.


## Running the application

Once the setup is complete, you can run the application in one of two ways:

**Method 1: Using `air` (Recommended for development)**

This method provides hot reload capabilities for faster development iterations.

```bash
air
```

This will start the application. You can access it via the localhost URL provided in the terminal output.


**Method 2: Using the standard `go run` command**

This method is suitable for simple execution, but lacks the hot-reloading feature of `air`.

```bash
go run *.go
```

This will compile and run the application.  You can then access it via the localhost URL indicated in the console output (likely port 8080, but check the output for confirmation).

## Demo Features Highlighted

This demo showcases the following Gemini Code Assist features:

* **Code Generation:**  Generating code snippets and entire functions based on natural language prompts.
* **Code Completion:**  Intelligently completing code based on context and existing code.
* **Code Explanation:**  Understanding and explaining existing code.
* **Code Refactoring:**  Improving code readability and maintainability.
* **Testing:** Generating unit tests for functions.
* **Security Analysis:** Identifying potential security vulnerabilities.

This application provides a practical example of how Gemini Code Assist can streamline the development process and improve code quality.
