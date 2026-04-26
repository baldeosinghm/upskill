# Phase 1

## Database Migrations

- Database: PostgreSQL
    - Reasoning: "Upskill" requires unique relationships with rigid data fields (i.e. email, courses, id) corresponing to an individual user.  This data must be quickly accesible as numerous values may be requested from multiple users at the same time. A relational database works best in this scenario.  Postgres is also great for data integrity, a must for a product that demands mutliple entities (i.e. students, teachers, courses) are receiving the appropriate work.
- Database Driver: `pgx`
    - pgx is a pure Go driver and toolkit for PostgreSQL designed for high performance and access to PostgreSQL-specific features.  It is widely considered the standard for Go/Postgres applications.

## Backend Architecture

- Used Handler-Service-Repository pattern
    - Advantages: decoupling, testing, flexibility.
    - Ensures each layer only executes functions in accordance to its purpose.

## HTTP Services
- Router: `chi`
    - An idiomatic router for building HTTP services in Go, `chi` is the ideal choice.

## Key Decisions
- Security
    - Used bcrypt.DefaultCost to ensure password integrity; the cost factor and salting increases the time required to hash each password, but makes brute-force attacks higher.
- Logging
    - Returned sentinel erorr, `ErrInvalidCredentials`, for both no-user and bad-password to prevent enumeration.
- Routes
    - Chose `/login` as opposed to more RESTful `/sessions` as endpoint for user login. `/sessions` reflects a strict REST/resource-oriented worldview ("a session is a thing you create"), while /login reflects an action-oriented worldview ("login is something you do"), which feels more in accordance with this product's current mental framework.

# Phase 2

## Key Decisions
- Enforcing Teacher-Owner constraint on courses table
    - There's two options here: enforce at the application layer or at the database level.  I chose to enforce this check at the service layer because
    if an attempt to make a non-teacher an owner of a course is made, it can be easily flagged.  Also, fixing this problem is easier here than at the DB-level.  If in the future, I add a role, it would involve appling a new constraint to the table and that won't as easy.