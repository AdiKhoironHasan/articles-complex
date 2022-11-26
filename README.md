# articles-complex

I use existing libs :

- Echo Framework
- Viper, for config management
- Sqlx, for database connection
- Redis, for caching

I'am build a service to manage data articles using golang, using postgresql for database, and redis for caching. service has database migration.
I created 5 endpoints :

1. [POST] /articles to post a new article
2. [GET] /articles to get the list of articles. Sort the articles by newest first with optional query parameters:
   - a. query: to search keywords in the article title and body
   - b. author: filter by author name.
3. [GET] /articles/:id to get a specific article by :id
4. [PUT] /articles/:id to update an article by :id
5. [DELETE] /articles/:id to delete an article by :id

# Setup after cloning the repo:

Run this command on your terminal to prepare dependencies:

- $ cd articles-complex
- $ go get
- $ go mod tidy

# connfiguration environment

Do this following actions for set up your configuration :

- copy config/config-dev.yaml.example ke config/config-dev.yaml
- complete the necessary credentials such as postgre and redis databases according to the existing format

# Database migration :

I use PostgreSQL for Database.
you can create database tables by migration, but before that you have to create a new database on your RDBMS.

- $ migrate -database "postgres://userDB:passwordDB@hostDB:portDB/yourDB?sslmode=disable" -path pkg/database/migrations up

# Run service :

You can run the service by using the following command, after that the service is ready to use.

- $ go run main.go

# Run with docker

You can run service in docker with this command:

- $ docker-compose -f docker-compose.yml up -d

# Use the service

You can deploy a service with consumption to an already created API endpoint, you can use the postman tool. To make it easier to use, I've created a workspace for it, and it's ready to go. don't forget to use postman with desktop agent if the service is running on localhost.
Postman Link:
https://api.postman.com/collections/18402968-801dd7d7-279b-424e-9c14-75fc4ac465e6?access_key=PMAT-01GJRV8GG1M6SQVNPRGD1F1RA7

# Thank You

If there are problems or want to know more information about me, please contact via linkedin via the following link https://www.linkedin.com/in/adi-khoiron-hasan or by sending an email to adikhoironhasan@gmail.com. Thank You :)
