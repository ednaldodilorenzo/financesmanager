package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;

@Entity
@Getter
public class Budget {

    @EmbeddedId
    private BudgetId id;

    @Column(nullable = false)
    private Integer year;

    @Column(nullable = false)
    private Long value;

    @ManyToOne
    @JoinColumns({
            // 'name' is the column in the transaction table
            // 'referencedColumnName' is the column name in the referenced table
            @JoinColumn(name = "category_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Category category;
}
