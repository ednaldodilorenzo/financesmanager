package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import lombok.AllArgsConstructor;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.io.Serializable;

@Embeddable
@Getter
@Setter
@EqualsAndHashCode
@AllArgsConstructor
@NoArgsConstructor
public class AccountId implements Serializable {

    @Column(name = "user_id", nullable = false)
    private Long userId;

    @Column(nullable = false)
    private Long id;
}
