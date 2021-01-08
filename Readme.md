# Golang BoilerPlater

## Migration

Migrations are like **version control for your database**, allowing your team to define and share the application's database schema definition. If you have ever had to tell a teammate to manually add a column to their local database schema after pulling in your changes from source control, you've faced the problem that database migrations solve.


## Execution

1. Run ```docker-compose up``` to spin up database and admire.
2. Open ```localhost:8080```,  select **system** to ```PostgreSQL``` and put username and password.
3. Build image using ```docker build -t golang-api .```
4. Run ```docker run golang-api``` to run the container.
