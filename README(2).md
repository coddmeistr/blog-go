
# Simple Blog server
This application tries to implement a simple blog backend


### Technology or what I learned
In this project I studied microservice architecture\
The folder structure of this application is very strange, because the resource part is written in Gin with MVC architecture, and the User microservice is written in Go-kit with a microservice approach\
I used PostgreSQL as the main database and GORM as the ORM\
Implemented JWT Authorization\
Learned to use Docker and DockerCompose


### How to use
Clone this repository, write your own PostgreSQL connection string and run the server






## API Reference

#### Get all users

```http
  GET /v1/user/
```
Get all users

#### Get one user by ID

```http
  GET /v1/user/:id
```
Gets one user with id in path

#### Login

```http
  PUT /v1/user/
```
Login

#### Delete user

```http
  DELETE /v1/user/:id
```
Delete user by id in path

#### Create user

```http
  POST /v1/user/
```
Create new user

#### Update contact info

```http
  PUT /v1/user/:id/contact
```
Update user's contact info

#### Update personal info

```http
  PUT /v1/user/:id/personal
```
Update user's personal info

#### Update personal info location

```http
  PUT /v1/user/:id/personal/location
```
Update user's personal info location

#### Create new like

```http
  POST /v1/res/like
```
Like some resource like article or comment


#### Create article

```http
  GET /v1/res/art
```
Create new article


#### Delete article

```http
  DELETE /v1/res/art/:id
```
Delete article by id in path


#### Get one article

```http
  GET /v1/res/art/:id
```
Get article by id in path

#### Get all articles

```http
  GET /v1/res/art/
```
Get all articles

#### Create comment

```http
  GET /v1/res/comm
```
Create new comment


#### Delete comment

```http
  DELETE /v1/res/comm/:id
```
Delete comment by id in path


#### Get one comment

```http
  GET /v1/res/comm/:id
```
Get comment by id in path

#### Get all comments

```http
  GET /v1/res/comm/
```
Get all comments

