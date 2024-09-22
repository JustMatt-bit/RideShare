# RideShare
Built for T120B165 module.

## What is it?
RideShare is a web-based application for users to efficiently share their rides with other people.
The web-application offers benefits not only to the user by reducing the cost of travel, but reduces traffic congestion and promotes environmental sustainability.

## Functional Requirements
### User Registration
- Users should be able to create an account by providing their personal information such as name, email, and password.
- Additionally, users should have the option to register using their existing external accounts (e.g., Google, Facebook) through OAuth2.

### User Authentication
- Users should be able to log in to their accounts using their email and password.
- Users who registered using OAuth2 should be able to log in using external accounts.

### User Profile
- Users should be able to view and update their profile information such as name, contact details, and preferences.
- Users should be able to fill out their vehicle information to be assignable for their rides.

### Ride Search & Booking
- Users should be able to search for available rides based on ride parameters (e.g. his starting location and required destination).
- Users should be able to book a ride by selecting a ride from the search results.

### Ride Management
- Users should be able to create a ride by providing details such as the starting location, destination, date, and time.
- Ride owners should be able to manage their rides by editing or deleting them.

### Pre-ride chat
- Users should be able to communicate with each other through a chat feature before the ride.

### User Feedback
- Users should be able to provide feedback about ride owners by giving the driver a score.

## Technology Stack
The RideShare application will be built using the following technology stack:

- Backend: Golang
- Frontend: Svelte
- Database: MySQL