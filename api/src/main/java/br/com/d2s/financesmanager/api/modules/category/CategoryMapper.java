package br.com.d2s.financesmanager.api.modules.category;

import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.model.Category;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.NullValuePropertyMappingStrategy;

@Mapper(componentModel = "spring",
        nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.SET_TO_NULL)
public interface CategoryMapper {

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "filter", ignore = true)
    Category toEntity(CategoryDto.Request dto);

    @Mapping(source = "id.id", target = "id")
    CategoryDto.Response toDto(Category category);
}
