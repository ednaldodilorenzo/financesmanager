package br.com.d2s.financesmanager.api.infra.ratelimiter;

import lombok.AllArgsConstructor;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;

import java.time.Duration;
import java.time.Instant;

@Component
@AllArgsConstructor
public class SlidingWindowRateLimiter {
    private final StringRedisTemplate stringRedisTemplate;
    public boolean isAllowed(String key, int maxRequests, Duration window) {
        var redisKey = "rate_limit: " + key;
        long now = Instant.now().toEpochMilli();
        long windowStart = now - window.toMillis();
        stringRedisTemplate.opsForZSet().removeRangeByScore(redisKey, 0, windowStart);
        var count = stringRedisTemplate.opsForZSet().zCard(redisKey);
        if (count != null && count >= maxRequests) {
            return false;
        }
        // Add current request
        stringRedisTemplate.opsForZSet().add(redisKey, String.valueOf(now), now);
        stringRedisTemplate.expire(redisKey, window.plusSeconds(1));
        return true;
    }
}
