# Use the migrate/migrate:latest image as the base
FROM migrate/migrate:latest

# Copy the wait-for-it.sh script into the container
COPY ./wait-for-it.sh /wait-for-it.sh

# Make the wait-for-it.sh script executable
RUN chmod +x /wait-for-it.sh

# Copy the migrations directory into the container
COPY ./migrations /migrations

# Set the entrypoint to use wait-for-it.sh by default
ENTRYPOINT ["sh", "/wait-for-it.sh"]

