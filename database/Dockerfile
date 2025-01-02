FROM debian:bullseye-slim

# Install necessary tools
RUN apt-get update 
   
RUN apt-get install -y \
    postgresql-client \
    wget
    
RUN rm -rf /var/lib/apt/lists/*

# Install migrate CLI
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar -xz -C /usr/local/bin \
    && chmod +x /usr/local/bin/migrate

# Set the working directory
WORKDIR /migrations

COPY ./migrations .

# Default command
CMD ["sh", "-c", "until pg_isready -h db -p 5432; do echo waiting for db; sleep 2; done; migrate -path /migrations -database $DATABASE_URL up"]

