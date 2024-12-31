# E-Commerce-Lite

Docker commands:

Check status of image:
docker-compose ps

Start container in detached mode:
docker-compose up -d

Shutdown container:
docker-compose down

Start up container from compose file:
docker-compose up -d

Setup .env file to configure the database credentials:
POSTGRES_PASSWORD=yourpassword
POSTGRES_USER=youruser
POSTGRES_DB=yourdb

To confirm the database is up and running connect to the instance:
psql -h localhost -U yourusername -d yourdatabase
Once in the shell type \l to see all available databases