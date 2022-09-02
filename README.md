# Role Play - Web Application Engineer: Portal

## Set environment variables

```
$ export PROJECT_ID=<Project ID>
$ export DATABASE_USER=<Database User>
$ export DATABASE_NAME=<Database Name>
$ export INSTANCE_CONNECTION_NAME=<Instance Connection Name>
$ export ADMIN_API_KEY=<Admin API Key>
$ export PORTAL_ENDPOINT=<Portal Endpoint>
```

## Run application locally

```
$ go run . 
```

## Initialize datasets

```
$ curl -d @image-hashes.json -H "Admin-API-Key: ${ADMIN_API_KEY}" -H "Content-Type: application/json" ${PORTAL_ENDPOINT}/admin/bulk/imagehashes
```

```
$ curl -X POST -H "Admin-API-Key: ${ADMIN_API_KEY}" -H "Content-Type: application/json" ${PORTAL_ENDPOINT}/admin/init/jobhistory
```

```
$ curl -X POST -H "Admin-API-Key: ${ADMIN_API_KEY}" -H "Content-Type: application/json" ${PORTAL_ENDPOINT}/admin/init/ranking
```