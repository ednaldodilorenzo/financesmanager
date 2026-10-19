package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.model.CategoryType;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

public class CategoryDto {

    public record Request(
            @NotBlank(message = "Category must be set.")
            String name,
            @NotNull(message = "Category type must be set.")
            CategoryType type){}

    public record Response(PublicId id, String name, CategoryType type){}
}
