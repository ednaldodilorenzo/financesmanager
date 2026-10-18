package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.MappedSuperclass;
import jakarta.persistence.PostLoad;
import jakarta.persistence.PostPersist;
import jakarta.persistence.Transient;
import org.springframework.data.domain.Persistable;

@MappedSuperclass
public abstract class AssignedIdEntity<ID> implements Persistable<ID> {

    // Starts as new; Hibernate flips it after loading or persisting, so save() calls persist()
    // instead of merge() and skips the existence SELECT.
    @Transient
    private boolean isNew = true;

    @Override
    public boolean isNew() {
        return isNew;
    }

    @PostLoad
    @PostPersist
    void markNotNew() {
        this.isNew = false;
    }
}

