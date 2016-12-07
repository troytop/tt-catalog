# stackato-mongo
Demo Stackato mongo app

# API 
This application has 2 API endpoints:

## GET /read
This just reads data from the database and returns it.

Example:
```
$ curl http://stackato-mongo.192.168.77.1.nip.io/read         
[
	{
		"id": "5801156b50bf89e3d14c3587",
		"user": "lozethghsj",
		"email": "lozethghsj@gmail.com",
		"password": "mvzloqkbvz"
	}
]
```

## POST /write
This will write an entry with random data into the database and return the data it inserted.

Example:
```
$ curl -X POST http://stackato-mongo.192.168.77.1.nip.io/write
{
	"id": "",
	"user": "lozethghsj",
	"email": "lozethghsj@gmail.com",
	"password": "mvzloqkbvz"
}
```

# Deploy
To deploy this application you can just do `cf push`

## Requirements
Requires HCF to have 1 service called `stackato-mongo`
This application will work with local mongo.

To make this service available follow these steps:

1. `hsm create-instance stackato.hpe.mongo 3.0`
  * enter the instance name as `stackato-mongo`
2. `cf hsm enable-service-instance stackato-mongo stackato-mongo`
  * make sure to add all the endpoints to the security group 
3. `cf create-service stackato-mongo default stackato-mongo`
