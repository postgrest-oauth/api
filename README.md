README
------

```
$ ./postgrest-oauth-server -h
2017/11/06 21:49:22 Started!
Usage of ./postgrest-oauth-server:
  -accessTokenJWTSecret string
    	Secret key for generating JWT access tokens (default "morethan32symbolssecretkey!!!!!!")
  -cookieBlockKey string
    	Block key for cookie creation. 16, 24 or 32 random symbols are valid (default "16charssecret!!!")
  -cookieHashKey string
    	Hash key for cookie creation. 64 random symbols recommended (default "supersecret")
  -dbConnString string
    	Database connection string (default "postgres://user:pass@localhost:5432/test?sslmode=disable")
  -refreshTokenJWTSecret string
    	Secret key for generating JWT refresh tokens (default "notlesshan32symbolssecretkey!!!!")
  -templateName string
    	Name of template html file (default "index.html")
  -templatePath string
    	Path to template html file. With trailing slash (default "./")
```