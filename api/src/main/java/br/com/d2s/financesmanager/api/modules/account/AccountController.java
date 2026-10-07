package br.com.d2s.financesmanager.api.modules.account;

import br.com.d2s.financesmanager.api.dto.AccountDto;
import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.modules.account.service.AccountService;
import jakarta.validation.Valid;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Objects;

@RestController
@RequestMapping("/accounts")
@AllArgsConstructor
public class AccountController {
    private AccountService accountService;

    @GetMapping("/")
    public List<AccountDto.Response> findAll(@AuthenticationPrincipal Jwt jwt) {
        return accountService.getAll(Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @GetMapping("/{id}")
    public AccountDto.Response findById(@PathVariable Long id, @AuthenticationPrincipal Jwt jwt) {
        return accountService.getOne(id, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public void post(@RequestBody @Valid AccountDto.Request dto, @AuthenticationPrincipal Jwt jwt) {
        accountService.create(dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PatchMapping("/{id}")
    public void patch(@PathVariable Long id, @AuthenticationPrincipal Jwt jwt,
                      @RequestBody @Valid AccountDto.Request dto) {
        accountService.update(id, dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }
}
