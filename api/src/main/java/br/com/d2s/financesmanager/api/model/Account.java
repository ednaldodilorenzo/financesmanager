package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;

@Entity
@Getter
@Setter
public class Account {
    @EmbeddedId
    private AccountId id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false)
    private String type;
}
