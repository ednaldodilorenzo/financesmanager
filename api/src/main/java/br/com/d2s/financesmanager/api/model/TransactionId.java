package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import lombok.Getter;

@Embeddable
@Getter
public class TransactionId {
    @Column(nullable = false, name = "user_id")
    private Long userId;

    @Column(nullable = false)
    private Long id;
}
