package br.com.d2s.financesmanager.api.modules.auth.service;

import br.com.d2s.financesmanager.api.modules.user.UserDao;
import lombok.AllArgsConstructor;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Component;

@Component
@AllArgsConstructor
public class UserDetailsServiceImpl implements UserDetailsService {

    private UserDao userDao;

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        var user = userDao.findByEmail(username).orElseThrow(() -> new UsernameNotFoundException(username));

        return User
                .withUsername(user.getId().toString())
                .password(user.getPassword())
                .roles("USER")
                .build();
    }
}
