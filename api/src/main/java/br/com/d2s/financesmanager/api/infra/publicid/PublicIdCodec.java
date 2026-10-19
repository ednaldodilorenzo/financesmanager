package br.com.d2s.financesmanager.api.infra.publicid;

import org.sqids.Sqids;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class PublicIdCodec {
    private static final int MIN_LENGTH = 8;

    private final Sqids sqids;

    // The alphabet must never change once ids are published: every link and stored id would stop decoding.
    public PublicIdCodec(@Value("${APP_PUBLIC_ID_ALPHABET}") String alphabet) {
        this.sqids = Sqids.builder().alphabet(alphabet).minLength(MIN_LENGTH).build();
    }

    public String encode(long id) {
        return sqids.encode(List.of(id));
    }

    /**
     * @throws IllegalArgumentException if the value is not a canonical encoding of a single id
     */
    public long decode(String value) {
        List<Long> numbers;
        try {
            numbers = sqids.decode(value);
        } catch (RuntimeException e) {
            throw new IllegalArgumentException("Invalid id: " + value, e);
        }
        // Sqids decodes several strings to the same number; accepting only the canonical one keeps ids unique.
        if (numbers.size() != 1 || !encode(numbers.getFirst()).equals(value)) {
            throw new IllegalArgumentException("Invalid id: " + value);
        }
        return numbers.getFirst();
    }
}
