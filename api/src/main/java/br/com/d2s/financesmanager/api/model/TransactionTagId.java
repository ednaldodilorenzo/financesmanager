package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import lombok.Getter;

@Embeddable
@Getter
public class TransactionTagId {
    @Column(nullable = false)
    private String tag;

    @Column(nullable = false, name = "user_id")
    private Long userId;

    @Column(nullable = false)
    private Long id;
}
