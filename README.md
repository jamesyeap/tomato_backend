# tomato_backend
A basic todo list for CVWO 2022 application!

## notes to self

To connect to postgresql db hosted on Heroku,
```
heroku pg:psql -a tomato-db
```

To specify the version of Go that Heroku should use in production, add this line to `go.mod`
```
// +heroku goVersion go1.17
```
This fixes the issue where the build is failing due to incompatibility issues where Go dependencies require a certain version of Go.