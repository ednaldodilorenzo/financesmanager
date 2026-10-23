package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Entity
@Getter
@Setter
@NoArgsConstructor
public class Account extends AssignedIdEntity<AccountId> {
    @EmbeddedId
    private AccountId id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false)
    private Character type;

    private Integer dueDay;

    public Account(AccountId id) {
        this.id = id;
    }
}
