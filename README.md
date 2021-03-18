# dashboardauth

Authentication system for a generic Internet of Things Dashboard.

 ## Using
Once the server is running, you can navigate to `http://your_ip_address:8080` or `https://your_ip_address:9090` in order to login.

The default email address is `e@g.c` and the default password is `test`

## Installation
### Using Docker

The simplest way to install this is to user Docker. With docker Installed, checkout docker-compose.yml and from the same directory, run:
```bash
docker-compose up
```
This will initialize 2 containers, one for the Go server and the other for Postgres.

### The hard way
1. Install Go, NPM, and Postgres to your machine.
2. Initialize Postgres with
```
user=postgres
password=postgres
dbName=iot_dashboard
```
3. Checkout all of the files to your `$GOPATH/src`.
4. Install all of the go dependencies. From the project root, run
```bash
go get ./...
```

5. To build the react.js code, `cd /iotdashboard/iotdbfrontend` and run
```bash
npm run build
```
6. To run the server, from the project root run
```bash
go run main.go
```

