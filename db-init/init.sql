CREATE DATABASE memo;
CREATE USER myuser WITH PASSWORD 'mypassword';
GRANT ALL PRIVILEGES ON DATABASE memo TO myuser;
ALTER DATABASE memo SET timezone TO 'UTC';