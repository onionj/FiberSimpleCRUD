
# Simple CRUD with Fiber, and Validator

## TODO:
- &#x2610; tests
- &#x2610; websocket
- &#x2610; real database
- &#x2610; redis for limiter middleware storage
- &#x2610; docker and docker compose

### Create:

```curl
curl --request POST \
  --url http://localhost:8080/api/v1/users \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "sa",
	"email": "onionj98@gmail.com",
	"job": {
		"type": "developer",
		"salary": 1
	}

}'
```

### Read:

```curl
curl --request GET \
  --url http://localhost:8080/api/v1/users/saman \
  --header 'Content-Type: application/json'
```

### Update:

```curl
curl --request PATCH \
  --url http://localhost:8080/api/v1/users \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "saman",
	"email": "onionj98@gmail.com",
	"job": {
		"type": "youtuber",
		"salary": 1
	}

}'
```

### Delete:

```curl
curl --request DELETE \
  --url http://localhost:8080/api/v1/users/saman
```