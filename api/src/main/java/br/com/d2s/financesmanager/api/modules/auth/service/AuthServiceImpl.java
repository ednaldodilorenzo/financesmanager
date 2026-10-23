package br.com.d2s.financesmanager.api.modules.auth.service;

import br.com.d2s.financesmanager.api.config.app.AppProperties;
import br.com.d2s.financesmanager.api.dto.SignupDto;
import br.com.d2s.financesmanager.api.infra.exception.ApiException;
import br.com.d2s.financesmanager.api.infra.exception.ApiExceptionType;
import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.dto.LoginDto;
import br.com.d2s.financesmanager.api.dto.RegisterDto;
import br.com.d2s.financesmanager.api.infra.queue.OutboxMessageEvent;
import br.com.d2s.financesmanager.api.infra.queue.QueueProducer;
import br.com.d2s.financesmanager.api.infra.token.JwtManager;
import br.com.d2s.financesmanager.api.model.User;
import br.com.d2s.financesmanager.api.modules.user.UserDao;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;

@Service
@AllArgsConstructor
@Slf4j
public class AuthServiceImpl implements AuthService {
    private final AuthenticationManager authenticationManager;
    private final JwtManager jwtManager;
    private final UserDao userDao;
    private final QueueProducer queueProducer;
    private final AppProperties appProperties;
    private PasswordEncoder passwordEncoder;

    @Override
    public LoginDto.Response login(LoginDto.Request loginReq) {
        Authentication auth =  authenticationManager.authenticate(new UsernamePasswordAuthenticationToken(loginReq.email(), loginReq.password()));

        String token = jwtManager.generateAccessToken(auth);
        return new LoginDto.Response(new PublicId(Long.parseLong(auth.getName())), "", token);
    }

    @Override
    public void startRegistrationProcess(RegisterDto dto) {
        if (userDao.existsByEmail(dto.email())) {
            throw new ApiException(ApiExceptionType.CONFLICT, "Email already exists");
        }

        String token = jwtManager.generateRegistrationToken(dto.email());
        var outboxEvent = OutboxMessageEvent
                .builder()
                .subject("Confirmação de email")
                .email(dto.email())
                .payload(("Clique no link abaixo para confirmar seu email<br><br><a href=\"%s/register/%s\">Clique " +
                        "aqui").formatted(appProperties.webUrl(), token))
                .build();
        queueProducer.sendToQueue("email:send", outboxEvent);
    }

    @Override
    @Transactional
    public void registerWithToken(SignupDto dto) {
        var email = jwtManager.extractRegistrationInfo(dto.token());
        var user = new User(null, dto.name(), email, passwordEncoder.encode(dto.password()), Instant.now(), Instant.now());
        userDao.save(user);
    }
}
