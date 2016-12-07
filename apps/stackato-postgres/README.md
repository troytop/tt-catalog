# stackato-postgres
Demo Stackato postgres app

# API 
This application has 2 API endpoints:

## GET /read
This just reads data from the database and returns it.

Example:
```
$ curl http://stackato-postgres.192.168.77.1.nip.io/read         
[
	{
		"ID": 1,
		"CreatedAt": "2016-10-14T15:14:20Z",
		"UpdatedAt": "2016-10-14T15:14:20Z",
		"DeletedAt": null,
		"user": "jymfcvqasi",
		"email": "jymfcvqasi@gmail.com",
		"password": "ew2ukz33z4"
	}
]
```

## POST /write
This will write an entry with random data into the database and return the data it inserted.

Example:
```
$ curl -X POST http://stackato-postgres.192.168.77.1.nip.io/write
{
	"ID": 2,
	"CreatedAt": "2016-10-14T15:25:38.577421739Z",
	"UpdatedAt": "2016-10-14T15:25:38.577421739Z",
	"DeletedAt": null,
	"user": "dmfafmkssg",
	"email": "dmfafmkssg@gmail.com",
	"password": "cr0v15juqq"
}
```

# Deploy
To deploy this application you can just do `cf push`

## Requirements
Requires HCF to have 1 service called `stackato-postgres`
This application will work with local postgres.

To make this service available follow these steps:

1. `hsm create-instance stackato.hpe.postgres 9.4`
  * enter the instance name as `stackato-postgres`
2. `cf hsm enable-service-instance stackato-postgres stackato-postgres`
  * make sure to add all the endpoints to the security group 
3. `cf create-service stackato-postgres default stackato-postgres`
