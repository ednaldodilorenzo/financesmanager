package br.com.d2s.financesmanager.api.infra.publicid;

/**
 * Database id as exposed by the API. Serialized as a Sqids string by {@link PublicIdConfig};
 * services and entities keep working with the raw {@code Long}.
 */
public record PublicId(long value) {
    public static PublicId of(Long value) {
        return value == null ? null : new PublicId(value);
    }
}
