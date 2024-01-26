Project Readme

a. Project Theme Name User Registration System

b. Brief Project Description The purpose of this project is to create a simple user registration system with a PostgreSQL database backend. The project allows users to register by providing a username and password, and the registration data is stored in a PostgreSQL database.

c. Temirkhan Alisher, Delegatuly Zhasulan, Tulkubaev Tair

e. Step-by-Step Instructions for Launching the Application

Set Up PostgreSQL Database:

Create a PostgreSQL database with the name "postgres."
Update the connection string in the const connStr variable inside the insertUser function with your PostgreSQL database credentials.
Run the Application:

Clone the repository.
Open a terminal and navigate to the project directory.
Run the following command to start the server:
go run main.go
The server will be running on http://localhost:8080.
Access the Registration Form:

Open a web browser and go to http://localhost:8080.
You should see the User Registration form.
Register a User:

Enter a username and password in the form.
Click the "Register" button.
The registration data will be stored in the PostgreSQL database.
Retrieve Users (Optional):

To retrieve a list of registered users, send a GET request to http://localhost:8080/users.
f. Tools Used, Links to Sources

Programming Language: Go (Golang)
Database: PostgreSQL
Web Framework: Standard Go net/http package
