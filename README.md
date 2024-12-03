# Explore GRPC Service

This is a grpc service written using go. This service provides below 4 grpc endpoints:

## Endpoint details

### ListLikedYou
This lists out all the users those have liked the recipient user.
- **Request format**: 
```
{
  "recipient_user_id": "1",
  "pagination_token": "0"
}
```
- **Response format**:
```
{"likers":[{"actor_id":"2","unix_timestamp":"1732997071"},{"actor_id":"3","unix_timestamp":"1733169871"}],"next_pagination_token":""}
```

### ListNewLikedYou
This lists all users who liked the recipient excluding those who have been liked in return
- **Request format**: 
```
{
  "recipient_user_id": "1",
  "pagination_token": "0"
}
```
- **Response format**:
```
{"likers":[{"actor_id":"2","unix_timestamp":"1732997071"}],"next_pagination_token":""}
```

### CountLikedYou
This counts the number of users who liked the recipient
- **Request format**: 
```
{
  "recipient_user_id": "1"
}
```
- **Response format**:
```
{"count":"2"}
```

### PutDecision
This records the decision of the actor to like or pass the recipient. This decision can be overwritten. If an entry doesn't exist, it would be inserted. If already exists, it would be updated.
- **Request format**: 
```
{
  "actor_user_id": "1",
  "recipient_user_id": "4",
  "liked_recipient": true
}
```
- **Response format**:
```
{"mutual_likes":false}
```

### Instructions
In order to **build, run and test** it locally, follow below steps: 
- Clone this repository
- Change directory to this cloned repository: `cd explore`
- Change directory to docker: `cd docker`
- Run this command to setup postgres instance and explore service: `docker compose up -d`

**Testing**
- Once the service is up and running, it would be accessible via `localhost:50051`.
- You can test it through postman, add a new grpc request and import explore-service.proto file.
- Now select an endpoint from explore-service grpc definition and pass request body as documented above, it should return the response as mentioned above because database is seeded with a few records for testing.