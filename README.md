
# Rebate Management System

A Go-based rebate management system for managing rebate programs, transactions, and claims. The system uses **AWS DynamoDB** for storage and **Gin** for its HTTP API framework. Includes optional caching, real-time tracking for claims, and Swagger documentation for API exploration.

---

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Ensure AWS Credentials Are Configured Locally](#ensure-aws-credentials-are-configured-locally)
- [Postman Collection](#postman-collection)
- [Local Deployment](#local-deployment)
- [API Endpoints](#api-endpoints)
- [Real-Time Dashboard for Claims](#real-time-dashboard-for-claims)
- [Swagger Documentation](#swagger-documentation)
- [Error Handling](#error-handling)
- [Example Usage](#example-usage)
   - [Create a Rebate Program](#create-a-rebate-program)
   - [Submit a Transaction](#submit-a-transaction)
   - [Calculate Rebate](#calculate-rebate)
   - [Submit a Claim](#submit-a-claim)
- [Troubleshooting](#troubleshooting)
   - [AWS Credentials Issues](#aws-credentials-issues)
   - [Network Connectivity](#network-connectivity)
   - [Check Time Synchronization](#check-time-synchronization)
   - [Manual Time Update](#manual-time-update)

---

## Features

1. **Rebate Management**:
   - Create rebate programs with specific eligibility criteria and validity periods.
   - Submit transactions tied to rebate programs.
   - Calculate rebates for transactions.

2. **Claim Tracking**:
   - Submit rebate claims and track their status.
   - Real-time tracking of approved, pending, and rejected claims.

3. **In-Memory Caching**:
   - Optimize frequently requested reports for better performance.

4. **Swagger Integration**:
   - Explore and test API endpoints with integrated Swagger documentation.

5. **Local Deployment**:
   - Easily deployable locally using Docker and Docker Compose.

6. **Scalability**:
   - Designed for scalable and efficient use in production environments.

---

## Getting Started

### Prerequisites

Ensure you have:
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)
- AWS credentials configured locally in `~/.aws/credentials` and `~/.aws/config` for accessing DynamoDB.

### Ensure AWS Credentials Are Configured Locally

1. Open your AWS credentials file located at `~/.aws/credentials`:
   ```plaintext
   [default]
   aws_access_key_id = YOUR_ACCESS_KEY
   aws_secret_access_key = YOUR_SECRET_KEY
   ```

2. Open your AWS config file located at `~/.aws/config`:
   ```plaintext
   [default]
   region = us-east-1
   ```

These files allow the application to authenticate with AWS services.

---

### Postman Collection

You can use the Postman collection provided in this repository to test the API endpoints easily.

1. Import the Postman collection file:
    - File Name: `rebate.postman_collection.json`
    - Location: Available in the `postman` directory of this repository.

2. Open Postman and navigate to the **Import** button.

3. Select the `rebate.postman_collection.json` file and import it into Postman.

4. Use the pre-configured requests to test all available API endpoints.

---

## Local Deployment

1. Clone this repository:
   ```bash
   git clone <repository_url>
   cd <repository_name>
   ```

2. Build and run the application:
   ```bash
   docker-compose up --build
   ```

3. Access the application:
   - API: `http://localhost:8080/api/`
   - Swagger Documentation: `http://localhost:8080/swagger/index.html`

4. Stop the application:
   ```bash
   docker-compose down
   ```

---

## API Endpoints

| Endpoint                          | Method | Description                             |
|-----------------------------------|--------|-----------------------------------------|
| `/api/rebate`                     | POST   | Create a new rebate program             |
| `/api/transaction`                | POST   | Submit a transaction                    |
| `/api/claim`                      | POST   | Submit a rebate claim                   |
| `/api/calculate`                  | GET    | Calculate rebate for a transaction      |
| `/api/reporting`                  | POST   | Generate a report                       |
| `/api/claims/status`              | GET    | Get the real-time claim status          |
| `/api/claim`                      | PUT    | Change the status of a claim            |

---

## Real-Time Dashboard for Claims

This feature provides real-time tracking of rebate claims, showing how many are approved, pending, or rejected.

1. Access the real-time dashboard:
   - URL: `http://localhost:8080/api/claims/status`

2. Response example:
   ```json
   {
     "approved": 10,
     "pending": 5,
     "rejected": 2
   }
   ```

Use this to monitor claim statuses dynamically.

---

## Swagger Documentation

Swagger is integrated to explore and test all API endpoints.

- URL: `http://localhost:8080/swagger/index.html`
- Navigate to this URL after starting the application to access the interactive documentation.

---

## Error Handling

The system implements robust error handling to ensure clear communication of issues. Below is a list of possible errors:

| Error Code                  | Status Code | Message                                            |
|-----------------------------|-------------|----------------------------------------------------|
| `err-internal-server-error` | 500         | An unexpected error occurred. Please try again later. |
| `err-invalid-interval`      | 400         | The provided date interval is invalid.            |
| `err-not-eligible`          | 400         | The user is not eligible for this action.         |
| `err-claim-not-found`       | 400         | No claim found with the specified ID.             |
| `err-rebate-already-claimed`| 400         | The rebate has already been claimed.              |
| `err-multiple-program-name` | 400         | A rebate program with the same name already exists.|
| `err-transaction-cannot-create`| 500      | The transaction could not be created.             |
| `err-transaction-not-found` | 400         | No transaction found with the specified ID.       |
| `err-rebate-not-found`      | 400         | No rebate program found with the specified ID.    |
| `err-failed-to-list-claims` | 500         | Failed to retrieve the list of claims.            |
| `err-failed-to-get-cache`   | 500         | Failed to fetch cached data.                      |
| `err-failed-to-store-cache` | 500         | Failed to store cache data.                       |

Each error is mapped to a clear status code and message for easier debugging.

---

## Example Usage

### Create a Rebate Program
```bash
curl -X POST http://localhost:8080/api/rebate -H "Content-Type: application/json" -d '{
  "ProgramName": "Holiday Sale",
  "Percentage": 15,
  "StartDate": "2025-01-01T00:00:00Z",
  "EndDate": "2025-01-31T23:59:59Z",
  "EligibilityCriteria": true
}'
```

### Submit a Transaction
```bash
curl -X POST http://localhost:8080/api/transaction -H "Content-Type: application/json" -d '{
  "Amount": 200.00,
  "Date": "2025-01-15T10:00:00Z",
  "RebateID": "uuid-of-rebate-program"
}'
```

### Calculate Rebate
```bash
curl -X GET http://localhost:8080/api/calculate?transactionId=<transactionId>
```

### Submit a Claim
```bash
curl -X POST http://localhost:8080/api/claim -H "Content-Type: application/json" -d '{
  "TransactionID": "uuid-of-transaction",
  "Date": "2025-01-15T10:00:00Z"
}'
```

---

## Troubleshooting

### AWS Credentials Issues
Ensure credentials are mounted correctly:
```bash
docker exec -it <container_name> cat /root/.aws/credentials
```

### Network Connectivity
Verify the container can reach AWS services:
```bash
docker exec -it <container_name> curl https://dynamodb.us-east-1.amazonaws.com
```

### Check Time Synchronization
Verify the container's time matches the host system:
```bash
docker exec -it <container_name> date
```

### Manual Time Update
If needed, manually synchronize the time in the container:
```bash
docker exec -it <container_name> ntpdate -u pool.ntp.org
```

---

Let me know if you need further changes!
