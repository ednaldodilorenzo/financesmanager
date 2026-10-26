package br.com.d2s.financesmanager.api.modules.auth;

import br.com.d2s.financesmanager.api.dto.LoginDto;
import br.com.d2s.financesmanager.api.dto.RegisterDto;
import br.com.d2s.financesmanager.api.dto.SignupDto;
import br.com.d2s.financesmanager.api.infra.security.TokenRevocationStore;
import br.com.d2s.financesmanager.api.modules.auth.service.AuthService;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.validation.Valid;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/auth")
@AllArgsConstructor
public class AuthController {

    private final AuthService authService;
    private final TokenRevocationStore tokenRevocationStore;

    @PostMapping("/login")
    public LoginDto.Response login(@RequestBody @Valid LoginDto.Request dto,
                                   HttpServletResponse resp) {
        LoginDto.Response response = authService.login(dto);
        var cookie = ResponseCookie
                .from("jwt", response.token())
                .httpOnly(Boolean.TRUE)
                .secure(Boolean.TRUE)
                .path("/")
                .maxAge(3600)
                .sameSite("Strict")
                .build();
        resp.addHeader(HttpHeaders.SET_COOKIE, cookie.toString());
        return response;
    }

    @PostMapping("/register")
    public void register(@RequestBody @Valid RegisterDto dto) {
        authService.startRegistrationProcess(dto);
    }

    @PostMapping("/signup")
    public void signup(@RequestBody @Valid SignupDto dto) {
        authService.registerWithToken(dto);
    }

    @PostMapping("/logout")
    public ResponseEntity<Void> logout(@AuthenticationPrincipal Jwt jwt) {
        if (jwt != null) {
            tokenRevocationStore.revoke(jwt.getId(), jwt.getExpiresAt());
        }
        ResponseCookie clear = ResponseCookie.from("jwt", "")
                .httpOnly(true).secure(true).path("/").maxAge(0).build();
        return ResponseEntity.noContent()
                .header(HttpHeaders.SET_COOKIE, clear.toString())
                .build();
    }
}
