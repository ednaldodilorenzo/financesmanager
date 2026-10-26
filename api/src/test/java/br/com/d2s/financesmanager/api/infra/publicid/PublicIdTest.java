package br.com.d2s.financesmanager.api.infra.publicid;

import br.com.d2s.financesmanager.api.config.app.AppProperties;
import org.junit.jupiter.api.Test;
import org.sqids.Sqids;
import tools.jackson.databind.DatabindException;
import tools.jackson.databind.json.JsonMapper;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class PublicIdTest {
    private static final String ALPHABET = "KIHZiM4R6o1p2hDOf8YXLeAJjw0UnPkgzdQaTGsNWmFyBcVqEr7uv5t3S9Clbx";
    AppProperties appProperties = new AppProperties("xx", "xx", "xxxxxxxx");
    private final PublicIdCodec codec = new PublicIdCodec(appProperties);
    private final PublicIdConfig config = new PublicIdConfig();
    private final JsonMapper mapper = JsonMapper.builder().addModule(config.publicIdModule(codec)).build();

    record Dto(PublicId id) {}

    @Test
    void encodesToOpaqueStringAndDecodesBack() {
        String encoded = codec.encode(15432L);

        assertTrue(encoded.length() >= 8);
        assertFalse(encoded.contains("15432"));
        assertEquals(15432L, codec.decode(encoded));
    }

    @Test
    void rejectsValuesThatAreNotASingleCanonicalId() {
        String twoNumbers = Sqids.builder().alphabet(ALPHABET).minLength(8).build().encode(List.of(1L, 2L));

        assertThrows(IllegalArgumentException.class, () -> codec.decode(twoNumbers));
        assertThrows(IllegalArgumentException.class, () -> codec.decode("15432"));
        assertThrows(IllegalArgumentException.class, () -> codec.decode("!!!!!!!!"));
    }

    @Test
    void jsonRoundTrip() {
        String json = mapper.writeValueAsString(new Dto(new PublicId(42L)));

        assertEquals("{\"id\":\"" + codec.encode(42L) + "\"}", json);
        assertEquals(new Dto(new PublicId(42L)), mapper.readValue(json, Dto.class));
    }

    @Test
    void jsonTreatsBlankAndNullAsMissing() {
        assertNull(mapper.readValue("{\"id\":\"\"}", Dto.class).id());
        assertNull(mapper.readValue("{\"id\":null}", Dto.class).id());
    }

    @Test
    void jsonRejectsRawNumericIds() {
        assertThrows(DatabindException.class, () -> mapper.readValue("{\"id\":\"42\"}", Dto.class));
        assertThrows(DatabindException.class, () -> mapper.readValue("{\"id\":42}", Dto.class));
    }

    @Test
    void converterDecodesPathVariables() {
        var converter = config.publicIdConverter(codec);

        assertEquals(new PublicId(7L), converter.convert(codec.encode(7L)));
        assertNull(converter.convert(""));
        assertThrows(IllegalArgumentException.class, () -> converter.convert("7"));
    }
}
