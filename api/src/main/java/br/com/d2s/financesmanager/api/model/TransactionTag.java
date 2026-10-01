package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.EmbeddedId;
import jakarta.persistence.Entity;
import lombok.Getter;

@Entity
@Getter
public class TransactionTag {

    @EmbeddedId
    private TransactionTagId id;
}
