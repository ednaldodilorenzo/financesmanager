# ==============================================================================
# STAGE 1: Build the Native Executable (GraalVM JDK 25 environment)
# ==============================================================================
FROM eclipse-temurin:25-jdk-jammy AS builder

# Set working directory inside the container
WORKDIR /app

# Install xargs (part of findutils) so the gradlew script can run
#RUN microdnf install -y findutils && microdnf clean all

# Copy Gradle wrapper and configuration files
COPY gradlew build.gradle.kts settings.gradle.kts ./
COPY gradle gradle

# Ensure gradlew has execution permissions
RUN chmod +x gradlew

# Download dependencies first
RUN ./gradlew dependencies --no-daemon

# Copy the actual source code
COPY src src

# Build the native image executable
#RUN ./gradlew nativeCompile --no-daemon
RUN ./gradlew bootJar -x test


# ==============================================================================
# STAGE 2: Run the Application (Minimal lightweight environment)
# ==============================================================================
FROM eclipse-temurin:25-jre-jammy

WORKDIR /app

# Copy the compiled native binary from the builder stage
# (Adjust 'my-app' to match your actual artifact name)
#COPY --from=builder /build/build/native/nativeCompile/api /app/server
COPY --from=builder /app/build/libs/*.jar app.jar

# Expose your Spring Boot port
EXPOSE 8080
ENV JAVA_OPTS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0 -Djava.security.egd=file:/dev/./urandom"

# Run the binary directly (no "java -jar" needed!)
#ENTRYPOINT ["/app/server"]
ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar app.jar"]
