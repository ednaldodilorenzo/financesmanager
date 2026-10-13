# ==============================================================================
# STAGE 1: Build the Native Executable (GraalVM JDK environment)
# STAGE 1: Build the Native Executable (GraalVM JDK 25 environment)
# ==============================================================================
FROM ghcr.io/graalvm/native-image-community:25 AS builder

# Set working directory inside the container
WORKDIR /build

# Copy Gradle wrapper and configuration files first (improves Docker layer caching)
COPY gradlew .
# Install xargs (part of findutils) so the gradlew script can run
RUN microdnf install -y findutils && microdnf clean all

# Copy Gradle wrapper and configuration files
COPY gradlew build.gradle.kts settings.gradle.kts ./
COPY gradle gradle

# Ensure gradlew has execution permissions
RUN chmod +x gradlew

# Download dependencies first (cached if build.gradle doesn't change)
# Download dependencies first
RUN ./gradlew dependencies --no-daemon

# Copy the actual source code
COPY src src

# Build the native image executable
RUN ./gradlew nativeCompile --no-daemon


# ==============================================================================
# STAGE 2: Run the Application (Minimal lightweight environment)
# ==============================================================================
FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

# Copy the compiled native binary from the builder stage
# (Adjust 'my-app' to match your actual artifact name)
COPY --from=builder /build/build/native/nativeCompile/api /app/server

# Expose your Spring Boot port
EXPOSE 8080

# Run the binary directly (no "java -jar" needed!)
ENTRYPOINT ["/app/server"]

