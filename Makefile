run-dev:
	docker-compose -f docker-compose-dev.yml up --build -d

run-prod:
        CERT_PATH=/home/edsf/tmp/certs MAIL_HOST=smtp.gmail.com MAIL_PORT=587 MAIL_USERNAME=ednado.dilorenzo@ifpb.edu.br MAIL_PASSWORD=123456 APP_URL=https://localhost JWT_KEY=xuxa docker-compose up --build -d

dev-logs:
	docker-compose -f docker-compose-dev.yml logs -f

down-dev:
	docker-compose -f docker-compose-dev.yml down
