# Buenavida API

## Setup Databases

This script generates the required project structure in MongoDB and PostgreSQL then exits

```
sh ./setup-db.sh
```

This one will run de databases indefinitely

```
sh ./start.sh
```

## Setup Go API

Move to `src` directory

```
cd src
```

Download dependencies

```
go mod tidy
```

Run

```
go run .
```
