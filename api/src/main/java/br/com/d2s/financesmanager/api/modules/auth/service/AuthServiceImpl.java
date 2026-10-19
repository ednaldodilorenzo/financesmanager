package br.com.d2s.financesmanager.api.modules.auth.service;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.dto.LoginDto;
import br.com.d2s.financesmanager.api.dto.RegisterDto;
import br.com.d2s.financesmanager.api.infra.token.JwtManager;
import lombok.AllArgsConstructor;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

@Service
@AllArgsConstructor
public class AuthServiceImpl implements AuthService {
    private final AuthenticationManager authenticationManager;
    private final JwtManager jwtManager;

    @Override
    public LoginDto.Response login(LoginDto.Request loginReq) {
        Authentication auth =  authenticationManager.authenticate(new UsernamePasswordAuthenticationToken(loginReq.email(), loginReq.password()));

        String token = jwtManager.generateAccessToken(auth);
        return new LoginDto.Response(new PublicId(Long.parseLong(auth.getName())), "", token);
    }

    @Override
    public void startRegistrationProcess(RegisterDto dto) {

    }
}
