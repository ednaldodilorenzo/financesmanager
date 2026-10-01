package br.com.d2s.financesmanager.api.modules.auth;

import br.com.d2s.financesmanager.api.dto.LoginDto;
import br.com.d2s.financesmanager.api.dto.RegisterDto;
import br.com.d2s.financesmanager.api.modules.auth.service.AuthService;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.validation.Valid;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseCookie;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/auth")
@AllArgsConstructor
public class AuthController {

    private final AuthService authService;

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
    public void register(@RequestBody RegisterDto dto) {

    }
}
