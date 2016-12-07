# stackato-redis
Demo Stackato redis app

# API 
This application has 2 API endpoints:

## GET /read
This just reads data from the database and returns it.

Example:
```
$ curl http://stackato-redis.192.168.77.1.nip.io/read         
[
	{
		"Key": "hblby",
		"Value": "c2237pi1px"
	}
]
```

## POST /write
This will write an entry with random data into the database and return the data it inserted.

Example:
```
$ curl -X POST http://stackato-redis.192.168.77.1.nip.io/write
{
	"Key": "hblby",
	"Value": "c2237pi1px"
}
```

# Deploy
To deploy this application you can just do `cf push`

## Requirements
Requires HCF to have 1 service called `stackato-redis`
This application will work with local redis.

To make this service available follow these steps:

1. `hsm create-instance stackato.hpe.redis 3.0`
  * enter the instance name as `stackato-redis`
2. `cf hsm enable-service-instance stackato-redis stackato-redis`
  * make sure to add all the endpoints to the security group 
3. `cf create-service stackato-redis default stackato-redis`
