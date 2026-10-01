package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import lombok.Getter;

@Embeddable
@Getter
public class CategoryId {

    private Long id;

    @Column(name = "user_id")
    private Long userId;
}
