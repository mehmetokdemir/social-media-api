## Social Media API
Rest API for social media app. Using golang, fiber, gorm, postgresql

Start golangci lint run
````shell
make lint
````

Make tidy for dependencies
````shell
make tidy
````

Build the database and app on docker
````shell
docker-compose build
````

Run with docker;
````shell
docker-compose up    
````

Down the docker;
````shell
docker-compose down    
````

Generate api docs;
````shell
make generate-docs
````

Swagger api documentation;
````shell
*/swagger/index.html
````
