package br.com.d2s.financesmanager.api.infra.token;

import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.jose.jws.MacAlgorithm;
import org.springframework.security.oauth2.jwt.JwsHeader;
import org.springframework.security.oauth2.jwt.JwtClaimsSet;
import org.springframework.security.oauth2.jwt.JwtEncoder;
import org.springframework.security.oauth2.jwt.JwtEncoderParameters;

import java.time.Instant;
import java.time.temporal.ChronoUnit;

public class JwtManagerImpl implements JwtManager {
    private final JwtEncoder encoder;

    public JwtManagerImpl(JwtEncoder encoder) {
        this.encoder = encoder;
    }

    @Override
    public String generateAccessToken(Authentication authentication) {
        Instant now = Instant.now();

        JwtClaimsSet claims = JwtClaimsSet.builder()
                .issuer("self")
                .issuedAt(now)
                .expiresAt(now.plus(1, ChronoUnit.HOURS))
                .subject(authentication.getName())
                .build();

        // Crucial step: Explicitly specify HS256 header parameters
        JwsHeader jwsHeader = JwsHeader.with(MacAlgorithm.HS256).build();

        return this.encoder.encode(JwtEncoderParameters.from(jwsHeader, claims)).getTokenValue();
    }
}
