package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Entity
@Getter
@Setter
@NoArgsConstructor
public class Category extends AssignedIdEntity<CategoryId> {
    @EmbeddedId
    private CategoryId id;

    @Column(nullable = false)
    private String name;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private CategoryType type;

    @Column(nullable = false)
    private String filter;

    public Category(CategoryId categoryId) {
        this.id = categoryId;
    }
}
