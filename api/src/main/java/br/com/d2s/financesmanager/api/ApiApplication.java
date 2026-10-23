package br.com.d2s.financesmanager.api;

import br.com.d2s.financesmanager.api.config.app.AppProperties;
import br.com.d2s.financesmanager.api.config.parser.csv.CsvProperties;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;

@SpringBootApplication
@EnableConfigurationProperties({CsvProperties.class, AppProperties.class})
public class ApiApplication {

    static void main(String[] args) {
        SpringApplication.run(ApiApplication.class, args);
    }

}
