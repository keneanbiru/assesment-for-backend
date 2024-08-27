postman link = https://group3-6966.postman.co/workspace/New-Team-Workspace~3641651c-2b3f-48fa-a55f-c2608c8610b7/collection/37367045-8d89f63a-56c3-422d-aa1e-bde60b6df4fa?action=share&creator=37367045

Loan Tracker API
Introduction

The Loan Tracker API is a Go-based backend service designed to manage user accounts, handle loan applications, and provide administrative functionalities. This API emphasizes secure authentication, role-based access control, and efficient management of user data and loan applications.
API Endpoints
User Functionalities
Retrieve User Profile

    Endpoint: GET /users/profile
    Description: Retrieve the authenticated user's profile.
    Response: Returns user profile data including ID, name, email, and timestamps.

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

