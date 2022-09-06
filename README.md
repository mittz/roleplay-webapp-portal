# Role Play - Web Application Engineer: Portal

## Initialize datasets

```
$ curl -d @image-hashes.json -H "Admin-API-Key: ${ADMIN_API_KEY}" -H "Content-Type: application/json" ${PORTAL_ENDPOINT}/admin/bulk/imagehashes
```

```
$ curl -d @users.json -H "Admin-API-Key: ${ADMIN_API_KEY}" -H "Content-Type: application/json" ${PORTAL_ENDPOINT}/admin/bulk/users
```

```
$ curl -X POST -H "Admin-API-Key: ${ADMIN_API_KEY}" -H "Content-Type: application/json" ${PORTAL_ENDPOINT}/admin/init/jobhistory
```

```
$ curl -X POST -H "Admin-API-Key: ${ADMIN_API_KEY}" -H "Content-Type: application/json" ${PORTAL_ENDPOINT}/admin/init/ranking
```