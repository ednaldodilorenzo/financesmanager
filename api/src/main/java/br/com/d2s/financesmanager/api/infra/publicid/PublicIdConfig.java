package br.com.d2s.financesmanager.api.infra.publicid;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.convert.converter.Converter;
import tools.jackson.core.JsonGenerator;
import tools.jackson.core.JsonParser;
import tools.jackson.databind.DeserializationContext;
import tools.jackson.databind.JacksonModule;
import tools.jackson.databind.SerializationContext;
import tools.jackson.databind.ValueDeserializer;
import tools.jackson.databind.ValueSerializer;
import tools.jackson.databind.module.SimpleModule;

@Configuration
public class PublicIdConfig {

    // JSON bodies (request and response). Spring Boot registers JacksonModule beans automatically.
    @Bean
    JacksonModule publicIdModule(PublicIdCodec codec) {
        return new SimpleModule("public-id")
                .addSerializer(PublicId.class, new ValueSerializer<>() {
                    @Override
                    public void serialize(PublicId id, JsonGenerator gen, SerializationContext ctxt) {
                        gen.writeString(codec.encode(id.value()));
                    }
                })
                .addDeserializer(PublicId.class, new ValueDeserializer<>() {
                    @Override
                    public PublicId deserialize(JsonParser p, DeserializationContext ctxt) {
                        String value = p.getValueAsString();
                        if (value == null || value.isBlank()) return null;
                        try {
                            return new PublicId(codec.decode(value));
                        } catch (IllegalArgumentException e) {
                            throw ctxt.weirdStringException(value, PublicId.class, "invalid id");
                        }
                    }
                });
    }

    // Path variables and form fields (@PathVariable, @ModelAttribute), which don't go through Jackson.
    @Bean
    Converter<String, PublicId> publicIdConverter(PublicIdCodec codec) {
        return source -> source.isBlank() ? null : new PublicId(codec.decode(source));
    }
}
