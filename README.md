# COWS API
A simple CRUD API to support the cow-names React UI. Allows for basic operations to store cow names.

![image](https://oldmooresalmanac.com/wp-content/uploads/2017/11/cow-2896329_960_720-Copy-476x459.jpg)

## API


### GET

**URL** : `/v0/cows/`

**Method** : `GET`

**Auth required** : NO

**Permissions required** : None

### Successful Response
**Code** : `200 OK`
**Content examples**

Returns an array of all cows on the farm.

```json
[
  {
    "name": "Curious",
    "id": "2927",
    "date": "12-03-2020",
    "image": "",
    "finder": "Alanya"
  },
  {
    "name": "Ferdinand",
    "id": "2766",
    "date": "7-03-2021",
    "image": "",
    "finder": "Daddy"
  }
]
```

-------
### GET by id
**URL** : `/v0/cows/{id}`

**Method** : `GET`

**Auth required** : NO

**Permissions required** : None

### Successful Response
**Code** : `200 OK`
**Content examples**

Returns an array with a single cow.

```json
[
  {
    "name": "Curious",
    "id": "2927",
    "date": "12-03-2020",
    "image": "",
    "finder": "Alanya"
  }
]
```
-------
### Create
**URL** : `/v0/cows/`

**Method** : `POST`

**Auth required** : NO

**Permissions required** : None

***Body***
```json
{
    "name": "Curious",
    "id": "2927",
    "date": "12-03-2020",
    "image": "",
    "finder": "Alanya"
  }
```

***Response***
```json
HTTP/1.1 200 OK
Vary: Origin
Date: Sun, 29 Aug 2021 15:22:02 GMT
Content-Length: 0
```

-------
### Update

**URL** : `/v0/cows/{id}`

**Method** : `PUT`

**Auth required** : NO

**Permissions required** : None

### Successful Response
**Code** : `200 OK`
**Content examples**

Updates a single cow's data.

***body***

```json
  {
    "name": "New Name",
    "id": "2927",
    "date": "12-03-2020",
    "image": "",
    "finder": "New Finder"
  }
```

***Response***
```json
HTTP/1.1 200 OK
Vary: Origin
Date: Sun, 29 Aug 2021 15:22:02 GMT
Content-Length: 0
```

-------
### Delete

**URL** : `/v0/cows/{id}`

**Method** : `DELETE`

**Auth required** : NO

**Permissions required** : None

### Successful Response
**Code** : `200 OK`
**Content examples**

Updates a single cow's data.

***body***

```json
  {
    "name": "New Name",
    "id": "2927",
    "date": "12-03-2020",
    "image": "",
    "finder": "New Finder"
  }
```

***Response***
```json
HTTP/1.1 200 OK
Vary: Origin
Date: Sun, 29 Aug 2021 15:22:02 GMT
Content-Length: 0
```



```
# runs tidy, format, lint, test
make check
```

# Build
From the root of this repository, run:
```shell
go build
```

# Run
Once built, the binaries will be under `target/`.
```shell
go run .
```