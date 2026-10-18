package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;

@Entity
@Getter
@Setter
public class Budget extends AssignedIdEntity<BudgetId> {

    @EmbeddedId
    private BudgetId id;

    @Column(nullable = false)
    private Integer year;

    @Column(nullable = false)
    private Long value;

    @Column(name = "category_id", nullable = false)
    private Long categoryId;

    @ManyToOne
    @JoinColumns({
            @JoinColumn(name = "category_id", referencedColumnName = "id", insertable = false, updatable = false),
            @JoinColumn(name = "user_id", referencedColumnName = "user_id", insertable = false, updatable = false)
    })
    private Category category;
}
