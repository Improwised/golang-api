# Golang BoilerPlater

### **Index**
- [Introduction](#introduction)

- [Multiple Database System Support](#multiple-database-system-support)

- [Kick Start Commands](#kick-start-commands)

- [Migrations](#migrations)

- [Code Walk-through](#code-walk-through)
    - [Config](#config)
    - [Command](#command)
    - [Route](#route)
    - [Middleware](#middleware)
    - [Model](#model)
    - [Controller](#controller)
    - [Utils](#utils)

- [Testcases](#testcases)
---
### **Introduction**
- This template is build on [gofiber](https://github.com/gofiber/fiber).
- To sync go packages you need to run following command.

    ### `go mod vendor`
- Copy `.env.example` as `.env` for environment variables.
- To add new package you can run following command.

    ### `go get {package_url}`
- You need to run `go mod vendor` after running above command to sync your pacakage with vendor directory.
- Following command is use to remove unused package from your `go.mod` and `go.sum` files.

    ### `go mod tidy`
- Please make sure you set proper database configuration to connect.
---
### **Multiple Database System Support**

- We are supporting 3 types of database at this time `postgres`, `mysql` & `sqlite3`

- It allows us to switch to database without changing too much stuff (but it requires to change function that are varies based on system.)

- [goqu](https://github.com/doug-martin/goqu) package allows us to do that.

- we have environment variable `DB_DIALECT` where you can set to database system type.
---
### **Kick Start Commands**
- This commands are define inside `Makefile`.

**Note:** Make sure you change this file if you are going to use this boilerplate with different folder structure.

for ex: `start` commands run `go run app.go api` if your app.go is somewhere else then you need to update this file for paths to work properly.

- **make start-dev** : This will start your app with `nodemon` which watch for changes in your directory and reload go.

- **make start** : To start api, it basically runs `go run app.go`

- **make create-migration file_name=migration_name** : This command will create migration `up` and `down` both respectively, you need to pass `file_name={migration_name}`

- **make build**  : It will build binary of your project.

- **make test** : To run testcase for entire project, `.env.testing` env need to use while writing testcase, you need to load 

- **make swagger-genrate** : To genrate swagger docs, This command will download binary if not exist.

- **make migration-up** : To run `Up` migrations.
- **make migration-down** : To run `Down` migrations, Which technically revert migration that are being perform by `Up`.
---
## Migration

Migrations are like **version control for your database**, allowing your team to define and share the application's database schema definition. If you have ever had to tell a teammate to manually add a column to their local database schema after pulling in your changes from source control, you've faced the problem that database migrations solve.


## Execution

1. Run ```docker-compose up``` to spin up database and admire.
2. Open ```localhost:8080```,  select **system** to ```PostgreSQL``` and put username and password.
3. Build image using ```docker build -t golang-api .```
4. Run ```docker run golang-api``` to run the container.

### **Migrations**
- **CREATE :** To create migrations i have discribed details in [Kick Start Commands](#kick-start-commands) section.
- **RUN :** To run migration there is two command `make migration-up` && `make migration-down`.
- Migration needs `-- +migrate Up` and `-- +migrate Down` respectively in starting of files, this is required because we are using [sql-migrate](https://github.com/rubenv/sql-migrate) package 



### **Code Walk-through**
- #### Config:
    - We are using [envconfig](https://github.com/kelseyhightower/envconfig) which binds env values to struct elements.
    ```go
        type Config struct {
            AppName string `envconfig:"APP_NAME"`
        }
    ```
    - In above example `APP_NAME` is bind with `AppName` element of struct so further when it is used it will have value from `env`.
    - To load env variables from `.env` we are using [gotdotenv](https://github.com/joho/godotenv), which basically initialized when app starts [main.go](app.go).
    - We have also create a seprate method `LoadTestEnv()` in [config](main.go) which helps to load envs from `.env.testing`
- #### Command:
    - We are using [cobra](https://github.com/spf13/cobra) for commands, you need create seprate file for your command.
    - To register your command add in [cli/main.go](cli/main.go)
    - After registering your command you can run that by following syntax `go run app.go {your_command_name}`
- #### Route:
    - If you are building api then you need to add routes in [routes/main.go](routes/main.go).
    - In that we have `Setup()` which will be initialized while starting app so all routes that are required need to register their.
    - Currently we are using seprate method to register routes, that sepration is based on controller. e.g If your controller have multiple or single we will create separte function to register it.
    ```go
        func setupAuthController(v1 fiber.Router, goqu *goqu.Database) error {
	        authController, err := controller.NewAuthController(goqu)
            if err != nil {
                return err
            }
            v1.Post("/login", authController.DoAuth)
            return nil
        }
    ```
    - Above function will be created inside [routes/main.go](routes/main.go), and that function will be called in `Setup()`.
    - Why we do this ? because every model does required database object that will perform database operation. We will be using same connection for all the database operation.
    - `controller.NewAuthController(goqu)` This line at starting of function calls controller method which calls initialize method of model.
    - So when we routes will register we will initialize method of controller, that will initialize model struct with database object.
- #### Middleware:
    - Currently we have two middleware [http_logger](middleware/http_logger.go) & [jwt](middleware/jwt.go)
    - [http_logger](middleware/http_logger.go): This middleware enusres that it logs every incoming request to logger.
    - [jwt](middleware/jwt.go): JWT validates the user authentication, in routes we need to set which routes need authentication, so middleware will check for auth on that request.
- #### Model:
    - Model refer to database tables, currently we have followed structure that one table have one model.
    - When you create a model it should inside of `models` folder.
    - Create a model struct with initialize method.
    ```go
    type UserModel struct {
	    db *goqu.Database
    }

    func InitUserModel(goqu *goqu.Database) (UserModel, error) {
        return UserModel{
            db: goqu,
        }, nil
    }
    ```
    - Above code snippet represent model struct and its init method that will be call in controller method to initialize database object.
    - When you create new method after initialization that function should be object of `UserModel` in current example. If you didn't do that then function will not able to access to database object.
    - Define struct that will represent your table columns.
    ```go
    type User struct {
        ID        string `json:"id"`
        FirstName string `json:"first_name" db:"first_name" validate:"required"`
        LastName  string `json:"last_name" db:"last_name" validate:"required"`
        Email     string `json:"email" db:"email" validate:"required"`
        Password  string `json:"password,omitempty" db:"password" validate:"required"`
    }
    ```
    - Above example represent fields that table have (we have added limited fields here)
    - `json` - in struct represent key name when it will be returned in json format.
    - `db` - represents name of database column, it's supported by `goqu` package.
    - `validate` - if you want to validate fields that are required you can set `validate:"required"`
    - `omitempty` - this will be used with along side of json, if value will be empty then that element will be trimmed from response of json. 
- #### Controller:
    - As descirbed in [route](#route) controller init method should be call in route.
    - Each controller must have their struct which contains model object, that will be use to call function of models.
- #### Utils:
    - We have define some common methods inside `utils` like json response.
    - We need common function which handles json response that's why we have created (json_response.go)[utils/json_response.go]
    - Similarly we have created different files for different use in utils.
---
### **Testcases**
- Test file should be available within the same package.
- It is advisible to create one `TestMain()` for one package.
- `TestMain()` should initialize everything you need in testcase.
- After that you can call `m.Run()` which executes all testcases that are available in testcases, Also `m.Run()` return  exit status code.
- After `m.RUN()` we must delete data that was inserted.
- When you run testcase with `make test` it will also display `code coverage` which helps you to determine percentage of code covered in your testcases.
