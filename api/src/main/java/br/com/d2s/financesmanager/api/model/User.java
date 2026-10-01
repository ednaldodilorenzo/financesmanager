package br.com.d2s.financesmanager.api.model;

import jakarta.persistence.*;
import lombok.Getter;

import java.time.Instant;

@Entity
@Getter
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String name;

    @Column(unique = true, nullable = false)
    private String email;

    @Column(nullable = false)
    private String password;

    @Column(nullable = false)
    private Instant createdAt;

    private Instant updatedAt;
}
