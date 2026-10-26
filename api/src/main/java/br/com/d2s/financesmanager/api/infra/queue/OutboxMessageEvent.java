package br.com.d2s.financesmanager.api.infra.queue;

import lombok.Builder;

@Builder
public record OutboxMessageEvent(
        String idempotencyKey,
        String email,
        String subject,
        String payload
) {
}
