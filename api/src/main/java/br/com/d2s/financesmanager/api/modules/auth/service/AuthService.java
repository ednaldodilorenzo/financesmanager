package br.com.d2s.financesmanager.api.modules.auth.service;

import br.com.d2s.financesmanager.api.dto.LoginDto;
import br.com.d2s.financesmanager.api.dto.RegisterDto;

public interface AuthService {
    LoginDto.Response login(LoginDto.Request loginReq);
    void startRegistrationProcess(RegisterDto dto);
}
