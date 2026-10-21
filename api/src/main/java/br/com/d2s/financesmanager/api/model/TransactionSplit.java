package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;

import java.util.HashSet;
import java.util.Set;


@Entity
@Getter
@Setter
public class TransactionSplit extends AssignedIdEntity<TransactionSplitId> {

    @EmbeddedId
    private TransactionSplitId id;

    @Column(nullable = false)
    private String description;

    private String detail;

    @Column(name = "category_id", nullable = false)
    private Long categoryId;

    @Column(name = "transaction_id", nullable = false)
    private Long transactionId;

    @ManyToOne
    @JoinColumns({
            @JoinColumn(name = "category_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Category category;

    @Column(nullable = false)
    private Long value;

    @ManyToOne(cascade = {CascadeType.ALL})
    @JoinColumns({
            @JoinColumn(name = "transaction_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Transaction transaction;

    @ElementCollection
    @CollectionTable(
            name = "transaction_tag",
            joinColumns = {
                    @JoinColumn(name = "transaction_split_id"),
                    @JoinColumn(name = "user_id")
            })
    @Column(name = "tag", nullable = false)
    private Set<String> tags = new HashSet<>();
}
