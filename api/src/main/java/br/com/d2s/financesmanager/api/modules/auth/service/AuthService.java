package br.com.d2s.financesmanager.api.modules.auth.service;

import br.com.d2s.financesmanager.api.dto.LoginDto;
import br.com.d2s.financesmanager.api.dto.RegisterDto;
import br.com.d2s.financesmanager.api.dto.SignupDto;

public interface AuthService {
    LoginDto.Response login(LoginDto.Request loginReq);
    void startRegistrationProcess(RegisterDto dto);
    void registerWithToken(SignupDto dto);
}
