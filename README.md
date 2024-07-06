# Spy Cat Agency

## Overview

The Spy Cat Agency Management System is a Go-based application that helps manage spy cats, their missions, and mission targets. This system allows for the creation, assignment, and tracking of missions for our elite team of feline operatives.

## Prerequisites

- Go
- PostgreSQL
- Echo web framework

## Setup

1. Clone the repository:
https://github.com/BloodsFa1zer/SpyCatAgency

2. Install dependencies:
`go mod download`

3. Set up the PostgreSQL database

4. Configure the database connection:
  _**".env"** file example can be seen in **".env.dist"** file_
  here is the copy:
   ```  
     POSTGRES_PASSWORD = user_password_to_database
     POSTGRES_USER = user_name_to_database
     POSTGRES_NAME = database_name
     POSTGRES_DRIVER = driver_database_name
     POSTGRES_PORT= database_port
   
    **That for Docker only:**

      DB_HOST= database_host
      POSTGRES_HOST= postgres_host

5. Run database migrations:
   ```  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest```
 
    and

    ``` migrate create -ext sql -dir database/migration/ -seq init_mg```

The application will start and listen on `http://localhost:6000` by default.


## API Endpoints

Here are some of the main API endpoints:

- `POST /cats` - Create a new Cat
- `GET /cats` - List all missions
- `GET /cats/:id` - Get details of a specific mission
- `PUT /cats` - Update a mission
- `DELETE /cats/:id` - Delete a mission
- `POST /missions` - Create a new mission
- `GET /missions` - List all missions
- `GET /missions/:id` - Get details of a specific mission
- `DELETE /missions/:id` - Delete a mission
- `POST /missions/:missionId/targets` - Add a target to a mission
- `DELETE /missions/:missionId/targets/:targetId` - Delete a target from a mission
- `PUT /missions/:missionId/targets/:targetId/complete` - Complete a target
- `PUT /missions/:missionId/assign` - Assign a cat to a mission
- `POST /missions/:missionId/targets` - Add a target to mission
- `PUT /targets` - Update a target
- `PUT /targets/:id/notes` - Update a target notes

