Here's the README for implementing the **Receipt Processor** web service using Go. This includes setup instructions, code explanations, and Docker setup:

---

# Receipt Processor (Go)

This is a web service written in Go that processes receipts and awards points based on specific rules. It includes two main endpoints: 
1. To process a receipt and generate an ID.
2. To retrieve the points awarded based on that receipt ID.

### Installation and Setup

#### Requirements
- Go 1.18+ installed
- Docker (optional, if you prefer Dockerized setup)

#### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Harshil-Gupta/receipt-processor.git
   cd receipt-processor
   ```

2. **Install dependencies**:
   Since this service uses Go’s built-in libraries for HTTP handling and UUID generation, there are no external dependencies.

3. **Build the application**:
   Run the following command to build the Go service:
   ```bash
   go build -o receipt-processor main.go
   ```

4. **Run the application**:
   After building the application, start the server using:
   ```bash
   ./receipt-processor
   ```
   The service will be available at `http://localhost:8080`.

#### Docker Setup (Optional)
1. **Build the Docker image**:
   Create a Docker image to containerize the application.
   ```bash
   docker build -t receipt-processor .
   ```

2. **Run the Docker container**:
   Start the Docker container to run the web service.
   ```bash
   docker run -p 8080:8080 receipt-processor
   ```
   The service should now be available at `http://localhost:8080`.

### API Documentation

#### 1. Process Receipt
   - **Endpoint**: `/receipts/process`
   - **Method**: `POST`
   - **Payload**: JSON object representing a receipt.
   - **Response**: JSON object with a unique ID for the receipt.

   **Example Payload**:
   ```json
   {
     "retailer": "Target",
     "purchaseDate": "2022-01-01",
     "purchaseTime": "13:01",
     "items": [
       { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
       { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
     ],
     "total": "35.35"
   }
   ```

   **Example Response**:
   ```json
   { "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
   ```

#### 2. Get Points
   - **Endpoint**: `/receipts/{id}/points`
   - **Method**: `GET`
   - **Response**: JSON object containing the points awarded for the receipt.

   **Example Response**:
   ```json
   { "points": 32 }
   ```

### Points Calculation Rules
The points for each receipt are calculated based on the following rules:

1. **Retailer Name**: 1 point for each alphanumeric character.
2. **Round Dollar Total**: 50 points if the total has no cents (e.g., `9.00`).
3. **Multiple of $0.25**: 25 points if the total is a multiple of `0.25`.
4. **Items**: 5 points for every two items on the receipt.
5. **Item Description Length**: For each item whose description length is a multiple of 3, multiply the item price by `0.2` and round up to the nearest integer. The result is added as points.
6. **Odd Purchase Day**: 6 points if the day is odd.
7. **Purchase Time**: 10 points if the purchase time is between 2:00 PM and 4:00 PM.

### Code Explanation

- **`main.go`**: This is the entry point of the service. It sets up the HTTP server, defines the routes, and handles the logic for both endpoints:
  - `/receipts/process`: Accepts receipt data, generates an ID, and stores the receipt in memory.
  - `/receipts/{id}/points`: Retrieves the stored receipt by its ID and calculates the points based on the defined rules.

- **In-Memory Storage**: The application uses a simple in-memory map to store receipt data temporarily. The data will be reset every time the application is restarted.

- **Helper Functions**:
  - **`calculatePoints`**: A function that calculates points for each receipt based on the rules.
  - **`generateID`**: A function that generates a unique ID for each processed receipt using the UUID library.
  - **`isOddDay`**, **`isBetweenTimes`**: Helper functions to check if the purchase day is odd or if the purchase time is between 2:00 PM and 4:00 PM.

### Running Tests
1. **Test API**: Use a tool like `curl` or Postman to test the service.
   - **POST /receipts/process**: Send the receipt JSON to this endpoint to get a unique ID.
     Example `curl` command:
     ```bash
     curl -X POST http://localhost:8080/receipts/process -H "Content-Type: application/json" -d @receipt.json
     ```
   - **GET /receipts/{id}/points**: Retrieve the points for the generated receipt ID.
     Example `curl` command:
     ```bash
     curl http://localhost:8080/receipts/{id}/points
     ```

### Example Walkthrough
1. **POST `/receipts/process`** with a sample receipt to generate an ID.
2. **GET `/receipts/{id}/points`** to retrieve the points for the generated receipt ID.

---

This README provides a comprehensive guide to building, running, and testing the **Receipt Processor** web service in Go. It also includes a Docker setup for easier deployment and testing.
