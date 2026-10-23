package br.com.d2s.financesmanager.api.infra.queue;

import lombok.AllArgsConstructor;
import org.springframework.data.redis.connection.stream.ObjectRecord;
import org.springframework.data.redis.connection.stream.StreamRecords;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Component;

@Component
@AllArgsConstructor
public class QueueProducer {
    private final RedisTemplate<String, Object> redisTemplate;

    public void sendToQueue(String queueName, OutboxMessageEvent event) {
        ObjectRecord<String, OutboxMessageEvent> record = StreamRecords.newRecord()
                .in(queueName)
                .ofObject(event);

        redisTemplate.opsForStream().add(record);
    }
}
