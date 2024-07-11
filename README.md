# Development environment installation

 1. Install PostgreSQL 16 and start it
 2. Create a user in PostgreSQL
 3. Create a database and make the previously created user an owner of the database
 4. Install Go lang 1.22.5+
 5. Run `go mod download`
 6. Run `go mod verfiy`
 7. Copy `example.env` config file and rename the copied file to `.env`
 8. Edit `.env` (port 8080 is preferred)
 9. Run `go run cmd/imi/college/main.go`
 10. Run the `create_database.sql` script on the created databse to fill the dictionaries data in
 11. The API now should be up and running
