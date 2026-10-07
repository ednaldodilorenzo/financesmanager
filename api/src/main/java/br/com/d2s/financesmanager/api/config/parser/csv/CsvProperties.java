package br.com.d2s.financesmanager.api.config.parser.csv;

import org.springframework.boot.context.properties.ConfigurationProperties;

import java.util.Map;

@ConfigurationProperties(prefix = "csv")
public record CsvProperties(Map<String, Layout> layouts) {
    public record Layout(int headerLines, int footerLines, int paymentDateColumn,
        int transactionDateColumn, int ammountColumn, String datePattern, String delimiter,
        int descriptionColumn){}
}
