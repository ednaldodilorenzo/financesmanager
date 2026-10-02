package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;

import java.time.Instant;

@Entity
@Getter
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

    @ManyToOne
    @JoinColumns({
            // 'name' is the column in the transaction table
            // 'referencedColumnName' is the column name in the referenced table
            @JoinColumn(name = "category_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Category category;

    @ManyToOne
    @JoinColumns({
            // 'name' is the column in the transaction table
            // 'referencedColumnName' is the column name in the referenced table
            @JoinColumn(name = "account_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Account account;
}
