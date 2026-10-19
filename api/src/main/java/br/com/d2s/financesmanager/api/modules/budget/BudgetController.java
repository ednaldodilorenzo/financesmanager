package br.com.d2s.financesmanager.api.modules.budget;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.dto.BudgetDto;
import br.com.d2s.financesmanager.api.modules.budget.service.BudgetService;
import jakarta.validation.Valid;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Objects;

@RestController
@RequestMapping("/v1/budgets")
@AllArgsConstructor
public class BudgetController {
    private BudgetService budgetService;

    @GetMapping("/")
    public List<BudgetDto.Response> getAll(@RequestParam Integer year, @AuthenticationPrincipal Jwt jwt) {
        return budgetService.findAll(year, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @GetMapping("/{id}")
    public BudgetDto.Response getById(@PathVariable PublicId id, @AuthenticationPrincipal Jwt jwt) {
        return budgetService.findById(id.value(), Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public void post(@RequestBody @Valid BudgetDto.Request budgetDto, @AuthenticationPrincipal Jwt jwt) {
        budgetService.create(budgetDto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PatchMapping("/{id}")
    public void patch(@PathVariable PublicId id, @RequestBody @Valid BudgetDto.Request budgetDto, @AuthenticationPrincipal Jwt jwt) {
        budgetService.update(id.value(), budgetDto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @DeleteMapping("/{id}")
    public void delete(@PathVariable PublicId id, @AuthenticationPrincipal Jwt jwt) {
        budgetService.delete(id.value(), Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }
}

