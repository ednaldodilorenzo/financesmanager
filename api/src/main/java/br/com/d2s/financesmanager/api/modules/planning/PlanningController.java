package br.com.d2s.financesmanager.api.modules.planning;

import br.com.d2s.financesmanager.api.dto.PlanningDto;
import br.com.d2s.financesmanager.api.modules.planning.service.PlanningService;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Objects;

@RestController
@RequestMapping("/v1/plannings")
@AllArgsConstructor
public class PlanningController {

    private final PlanningService planningService;

    @GetMapping("/")
    public List<PlanningDto> findPlanningByMonthAndYearAndUserId(@RequestParam Integer month, Integer year,
                                                                 @AuthenticationPrincipal Jwt jwt) {
        return planningService.findPlanningByMonthAndYear(month, year,
                Long.parseLong((Objects.requireNonNull(jwt.getSubject()))));
    }
}
