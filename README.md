## This is Example RESTful API server

**Project build in accordance with the concept of "Clean Architecture" described in this article** https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
REST API with custom JWT-based authentication system. Core functionality is about creating and managing storage of photos.

Structure:

4 Domain layers:

    Models layer
    Repository layer
    UseCase layer
    Delivery layer
    
## API:

### POST /auth/sign-up

Creates new user 

##### Example Input: 
```
{
	"username": "user",
	"password": "mysecretpassword"
} 
```


### POST /auth/sign-in

Request to get JWT Token based on user credentials

##### Example Input: 
```
Request Headers:
Authorization: "Basic dXNlcjU6bXlzZWN1cmVwYXNzd29yZA=="
```

##### Example Response: 
```
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzEwMzgyMjQuNzQ0MzI0MiwidXNlciI6eyJJRCI6IjAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMCIsIlVzZXJuYW1lIjoiemhhc2hrZXZ5Y2giLCJQYXNzd29yZCI6IjQyODYwMTc5ZmFiMTQ2YzZiZDAyNjlkMDViZTM0ZWNmYmY5Zjk3YjUifX0.3dsyKJQ-HZJxdvBMui0Mzgw6yb6If9aB8imGhxMOjsk"
} 
```
### DELETE /delete/user

Delete current user

#### Example Input:
```
{
  "delete":"ok"
}
```
##### Example Response:
```
{
  "response": "delete success"
}
```

### POST /api/photogramm/upload

Upload new photo

##### Example Input: 
```
Request Headers:
Content-Type: multipart/form-data; "boundary=--------------------------867303760029570575153177"
Request Body:
photo: undefined
```
##### Example Response:
```
{
    "responce": "20072ab149d6c20d96f73bffd4bf1628fc427de1"
}
```
### GET /api/photogramm/getphoto

Returns info about photo by id

#### Example Input:
```
{
    "id":"20c9f0501316ae4a8b130b0874a828d34f0e0252"
}
```
##### Example Response:
```
{
    "photo_id": "20c9f0501316ae4a8b130b0874a828d34f0e0252",
    "username": "user5",
    "user_id": 20,
    "album_name": "summer"
}
```
### DELETE /api/photogramm/removephoto

Remove photo by id

#### Example Input:
```
{
    "id":"20c9f0501316ae4a8b130b0874a828d34f0e0252"
}
```
##### Example Response:
```
{
  "response": "delete success"
}
```
### POST /api/photogramm/createalbum

Creates new photo album

#### Example Input:
```
{
    "name":"my new album"
}
```
##### Example Response:
```
{
  "response": "success"
}
```

### GET    /api/photogramm/getalbum

Returns info about album by name

#### Example Input:
```
{
    "name":"summer"
}
```
##### Example Response:
```
{
    "album_name": "summer",
    "photos": [
        "0aba003e44f59d2bf550edb65ce5f8483eb0947e",
        "20c9f0501316ae4a8b130b0874a828d34f0e0252",
        "c832e3a20d1d4454c7539ef27fb0c2f81a27c0ca",
        "cb4d0853d31bbea86e9fa17b40033f368104ce09"
    ],
    "total": 4
}
```
### DELETE /api/photogramm/removealbum

Remove album and photos belong to it 

#### Example Input:
```
{
    "name":"my new album"
}
```
##### Example Response:
```
{
  "response": "success"
}
```
### GET /api/photogramm/getinfo

Returns info about all users albums and photos

##### Example Response: 
```
{
    "id": 20,
    "username": "user5",
    "photo_albums": [
        {
            "album_name": "",
            "photos": [
                "41988035482768c411b75c83601946299ca1debf",
                "5afa7809707e20de224e1acb643f8b7fd70accb3",
                "5f90318c1c23d18bceacb1c15a24feeb9ae48162"
            ],
            "total": 3
        },
        {
            "album_name": "summer",
            "photos": [
                "0aba003e44f59d2bf550edb65ce5f8483eb0947e",
                "20c9f0501316ae4a8b130b0874a828d34f0e0252",
                "c832e3a20d1d4454c7539ef27fb0c2f81a27c0ca",
                "cb4d0853d31bbea86e9fa17b40033f368104ce09"
            ],
            "total": 4
        }
    ]
}
```

## Requirements
- go 1.19
- docker & docker-compose
