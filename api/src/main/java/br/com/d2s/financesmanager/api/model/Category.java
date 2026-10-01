package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;

@Entity
@Getter
public class Category {
    @EmbeddedId
    private CategoryId categoryId;

    @Column(nullable = false)
    private String name;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private CategoryType type;
}
