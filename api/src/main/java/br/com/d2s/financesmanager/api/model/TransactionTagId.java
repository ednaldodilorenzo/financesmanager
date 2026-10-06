package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import lombok.EqualsAndHashCode;
import lombok.Getter;

import java.io.Serializable;

@Embeddable
@Getter
@EqualsAndHashCode
public class TransactionTagId implements Serializable {
    @Column(nullable = false)
    private String tag;

    @Column(nullable = false, name = "user_id")
    private Long userId;

    @Column(name = "transaction_id", nullable = false)
    private Long transactionId;
}
