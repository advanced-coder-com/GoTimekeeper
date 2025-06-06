# Timekeeper GoLand project

## How to install / run

### 1. Install docker and go locally

### 2. Create `.env` file in this folder with content 

```
DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=timekeeper
APP_PORT=8080
DEBUG_PORT=2345
JWT_SECRET=supersecretkey
```

Change `JWT_SECRET`, `DB_USER`, `DB_PASSWORD` values

### 3. Build docker with `docker compose build`
### 4. Run project with `docker compose up`

## How to debug

Create Go Remote config with host `localhost` and port `2345`

