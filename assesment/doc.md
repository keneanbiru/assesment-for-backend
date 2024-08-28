Project Documentation for Loan Tracker API
Table of Contents

    Introduction
    API Endpoints
        User Functionalities
        Admin Functionalities
        Loan Management
        System Logs
    Non-Functional Requirements
    Setup and Running the Project
    Testing
    Postman Documentation
    Source Code
    Submission

Introduction

The Loan Tracker API is a backend service implemented in Go designed to manage user accounts, handle loan applications, and provide administrative functionalities. This API emphasizes secure authentication, role-based access control, and efficient management of user data and loan applications.
API Endpoints
User Functionalities

    Retrieve User Profile
        Endpoint: GET /users/profile
        Description: Retrieve the authenticated user's profile.
        Response: Contains user profile data including ID, name, email, and timestamps.

    Send Password Reset Link
        Endpoint: POST /users/password-reset
        Description: Send a password reset link to the user's email.
        Response: Indicates success or failure of the password reset request.

    Update Password After Reset
        Endpoint: POST /users/password-update
        Description: Update the user's password using the token received in the password reset email.
        Response: Confirms whether the password update was successful or not.

Admin Functionalities

    View All Users
        Endpoint: GET /admin/users
        Description: Retrieve a list of all users.
        Response: Provides a list of users with details including ID, name, email, and timestamps.

    Delete User Account
        Endpoint: DELETE /admin/users/{id}
        Description: Delete a specific user account.
        Response: Indicates success or failure of the delete operation.

Loan Management

    Apply for Loan
        Endpoint: POST /loans
        Description: Submit a loan application.
        Response: Provides the status of the loan application.

    View Loan Status
        Endpoint: GET /loans/{id}
        Description: Retrieve the status of a specific loan.
        Response: Details the loan status including ID, amount, status, and timestamps.

    View All Loans (Admin)
        Endpoint: GET /admin/loans
        Description: Retrieve all loan applications.
        Parameters: Allows filtering by loan status and sorting by order.
        Response: Provides a list of loan applications with details.

    Approve/Reject Loan (Admin)
        Endpoint: PATCH /admin/loans/{id}/status
        Description: Approve or reject a loan application.
        Response: Confirms the updated status of the loan.

    Delete Loan (Admin)
        Endpoint: DELETE /admin/loans/{id}
        Description: Delete a specific loan application.
        Response: Indicates success or failure of the delete operation.

System Logs

    View System Logs
        Endpoint: GET /admin/logs
        Description: Retrieve system logs.
        Response: Provides details of system logs including timestamps, events, and details.

Non-Functional Requirements

    Security
        Secure Password Handling: Use bcrypt for secure password hashing.
        JWT Authentication: Implement access and refresh tokens for stateless authentication.
        Role-Based Access Control (RBAC): Ensure proper authorization for admin functionalities.

    Performance
        Concurrency with Goroutines: Utilize goroutines to handle multiple tasks concurrently.
        API Response Optimization: Implement pagination and caching where necessary.

    Documentation
        Postman Documentation: Provide comprehensive documentation of all API endpoints, including example requests and responses.

Setup and Running the Project

    Clone the Repository
        Clone the GitHub repository to your local machine.

    Install Dependencies
        Download and install the necessary Go dependencies.

    Run the Application
        Start the application by running the main application file.

    Environment Variables
        Configure environment variables as required (e.g., database connection strings).

Testing

    Unit Tests
        Implement and run unit tests for key functionalities such as loan processing and user management.

    Run Tests
        Execute tests to ensure all functionalities work as expected.

Postman Documentation

    Link to Postman Collection: Postman Documentation