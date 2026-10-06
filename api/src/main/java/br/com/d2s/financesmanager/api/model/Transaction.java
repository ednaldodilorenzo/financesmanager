package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;

import java.time.Instant;

@Entity
@Getter
@Setter
public class Transaction {
    @EmbeddedId
    private TransactionId id;

    @Column(nullable = false)
    private String description;

    @Column(nullable = false)
    private Long value;

    @Column(nullable = false)
    private Instant paymentDate;

    @Column(nullable = false)
    private Instant transactionDate;

    private String detail;

    // FK columns are written through these basic fields; user_id is written by the embedded id.
    @Column(name = "category_id", nullable = false)
    private Long categoryId;

    @Column(name = "account_id", nullable = false)
    private Long accountId;

    // Read-only associations: Hibernate requires every column of a composite join
    // to be either insertable or not, and user_id is already owned by TransactionId.
    @ManyToOne
    @JoinColumns({
            @JoinColumn(name = "category_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Category category;

    @ManyToOne
    @JoinColumns({
            @JoinColumn(name = "account_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Account account;
}
