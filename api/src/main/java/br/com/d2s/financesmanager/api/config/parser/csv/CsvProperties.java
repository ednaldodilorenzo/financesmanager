package br.com.d2s.financesmanager.api.config.parser.csv;

import org.springframework.boot.context.properties.ConfigurationProperties;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.util.Map;

@ConfigurationProperties(prefix = "csv")
public record CsvProperties(Map<String, Layout> layouts) {
    public record Layout(int headerLines, int footerLines, int paymentDateColumn,
        int transactionDateColumn, int ammountColumn, String datePattern, String delimiter,
        int descriptionColumn, Charset charset){
        public Layout {
            if (charset == null) {
                charset = StandardCharsets.ISO_8859_1;
            }
        }
    }
}
