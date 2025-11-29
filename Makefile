generate_mapping:
	protoc -I ./mapping \
			--go_out=. \
			--go-grpc_out=. \
			./mapping/api/mapping.proto

generate_tokenizer:
	protoc -I ./tokenizer \
			--go_out=. \
			--go-grpc_out=. \
			./tokenizer/api/tokenizer.proto

generate_auth:
	protoc -I ./auth_service \
			--go_out=. \
			--go-grpc_out=. \
			./auth_service/api/auth_service.proto

generate_certs:
	openssl req -newkey rsa:2048 \
      -nodes -x509 \
      -days 3650 \
      -out certs/ca.pem \
      -keyout certs/ca.key \
      -subj "/C=RU/ST=VologdaOblast/L=Vologda/O=anonix/OU=dev/CN=localhost"

      