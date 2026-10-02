package br.com.d2s.financesmanager.api.modules.category;

import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.model.Category;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface CategoryMapper {

    Category toEntity(CategoryDto.Request dto);

    @Mapping(source = "categoryId.id", target = "id")
    CategoryDto.Response toDto(Category category);
}
