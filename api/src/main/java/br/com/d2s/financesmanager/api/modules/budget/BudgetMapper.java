package br.com.d2s.financesmanager.api.modules.budget;

import br.com.d2s.financesmanager.api.infra.publicid.PublicIdMapper;
import br.com.d2s.financesmanager.api.dto.BudgetDto;
import br.com.d2s.financesmanager.api.model.Budget;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;


@Mapper(componentModel = "spring", uses = PublicIdMapper.class)
public interface BudgetMapper {

    Budget toEntity(BudgetDto.Request budgetDto);

    @Mapping(source = "id.id", target = "id")
    BudgetDto.Response toDto(Budget budget);
}
