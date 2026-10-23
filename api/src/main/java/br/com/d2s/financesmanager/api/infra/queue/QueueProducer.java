package br.com.d2s.financesmanager.api.infra.queue;

import lombok.AllArgsConstructor;
import org.springframework.data.redis.connection.stream.MapRecord;
import org.springframework.data.redis.connection.stream.StreamRecords;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;

import java.util.Map;

@Component
@AllArgsConstructor
public class QueueProducer {
    private final StringRedisTemplate redisTemplate;

    public void sendToQueue(String queueName, OutboxMessageEvent event) {
        MapRecord<String, String, String> record = StreamRecords.newRecord()
                .in(queueName)
                .ofMap(Map.of(
                        "to", event.email(),
                        "subject", event.subject(),
                        "body", event.payload()));

        redisTemplate.opsForStream().add(record);
    }
}
