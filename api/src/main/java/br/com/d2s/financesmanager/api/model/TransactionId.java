package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import lombok.AllArgsConstructor;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;

import java.io.Serializable;

@Embeddable
@Getter
@EqualsAndHashCode
@AllArgsConstructor
@NoArgsConstructor
public class TransactionId implements Serializable {
    @Column(nullable = false, name = "user_id")
    private Long userId;

    // @GeneratedValue is not applied inside an @EmbeddedId; the value comes from
    // transaction_id_seq (see TransactionDao#nextId).
    @Column(nullable = false)
    private Long id;
}
