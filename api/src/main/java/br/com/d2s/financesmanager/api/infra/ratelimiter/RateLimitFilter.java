package br.com.d2s.financesmanager.api.infra.ratelimiter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.AllArgsConstructor;
import org.jspecify.annotations.NonNull;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.time.Duration;

@Component
@Order(1)
@AllArgsConstructor
public class RateLimitFilter extends OncePerRequestFilter {
    private final SlidingWindowRateLimiter rateLimiter;

    @Override
    protected void doFilterInternal(@NonNull HttpServletRequest request,
                                    @NonNull HttpServletResponse response,
                                    @NonNull FilterChain chain) throws ServletException, IOException {
        String clientIp = getClientIp(request);
        String endpoint = request.getRequestURI();
        // Stricter limits for auth endpoints
        int maxRequests = endpoint.startsWith("/api/auth") ? 5 : 100;
        Duration window = endpoint.startsWith("/api/auth")
                ? Duration.ofMinutes(1) : Duration.ofSeconds(10);
        if (!rateLimiter.isAllowed(clientIp + ":" + endpoint, maxRequests, window)) {
            response.setStatus(429);
            response.setHeader("Retry-After", "60");
            response.getWriter().write("{\"error\":\"Rate limit exceeded\"}");
            return;
        }
        chain.doFilter(request, response);
    }
    private String getClientIp(HttpServletRequest request) {
        String xForwardedFor = request.getHeader("X-Forwarded-For");
        if (xForwardedFor != null && !xForwardedFor.isEmpty()) {
            return xForwardedFor.split(",")[0].trim();
        }
        return request.getRemoteAddr();
    }
}
