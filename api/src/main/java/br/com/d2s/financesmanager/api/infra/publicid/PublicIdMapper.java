package br.com.d2s.financesmanager.api.infra.publicid;

import org.springframework.stereotype.Component;

// Lets MapStruct convert between entity Long ids and DTO PublicIds (referenced via @Mapper(uses = ...)).
@Component
public class PublicIdMapper {
    public PublicId toPublicId(Long id) {
        return PublicId.of(id);
    }

    public Long toLong(PublicId id) {
        return id == null ? null : id.value();
    }
}
