VERSION?=latest

docker/image/enum-batch-gen:
	docker build -t enum-batch-gen-build -f ./enum-batch-gen/Dockerfile.build .
	docker create --name enum-batch-gen-build enum-batch-gen-build
	docker cp enum-batch-gen-build:/go/bin/enum-batch-gen ./build/enum-batch-gen
	docker rm enum-batch-gen-build
	docker build -t enum-batch-gen -f ./enum-batch-gen/Dockerfile .

docker/image/enum:
	docker build -t key-crawler-enum -f enum/Dockerfile ./enum

docker/image/queuer:
	docker build -t queuer-build -f ./queuer/Dockerfile.build .
	docker create --name queuer-build queuer-build
	docker cp queuer-build:/go/bin/queuer ./build/
	docker rm queuer-build
	docker build -t key-crawler-queuer -f ./queuer/Dockerfile .

docker/image/recorder:
	docker build -t recorder-build -f ./recorder/Dockerfile.build .
	docker create --name recorder-build recorder-build
	docker cp recorder-build:/go/bin/recorder ./build/
	docker rm recorder-build
	docker build -t key-crawler-recorder -f ./recorder/Dockerfile .

docker/image/all: docker/image/enum docker/image/queuer docker/image/recorder docker/image/enum-batch-gen

tf/apply/dev:
	terraform apply -var-file=infra/dev.tfvars infra/

tf/destroy/dev:
	terraform destroy -var-file=infra/dev.tfvars infra/

compose/up:
	BATCH_QUEUE_URL="$$(terraform output enum_batch_queue_url)" \
		KEY_QUEUE_URL="$$(terraform output keys_queue_url)" \
		docker-compose up --build

aws/ecr/signin:
	`aws ecr get-login --no-include-email --region us-east-1`

# -------------------
#
#images/enum-enum-batch-gen: build/enum-enum-batch-gen
#	docker build -t enum-enum-batch-gen -f ./enum-enum-batch-gen/Dockerfile .
#
#build/enum-enum-batch-gen:
#	docker build -t enum-enum-batch-gen-build -f ./enum-enum-batch-gen/Dockerfile.build .
#	docker create --name enum-enum-batch-gen-build enum-enum-batch-gen-build
#	docker cp enum-enum-batch-gen-build:/go/bin/enum-batch-gen ./build/enum-enum-batch-gen
#
#images/recorder: build/recorder
#	docker build -t key-crawler-recorder -f ./recorder/Dockerfile .
#
#build/recorder:
#	docker build -t recorder-build -f ./recorder/Dockerfile.build .
#	docker create --name recorder-build recorder-build
#	docker cp recorder-build:/go/bin/recorder ./build/
#
#images/queuer: build/queuer
#	docker build -t key-crawler-queuer -f ./queuer/Dockerfile .
#
#build/queuer:
#	docker build -t queuer-build -f ./queuer/Dockerfile.build .
#	docker create --name queuer-build queuer-build
#	docker cp queuer-build:/go/bin/queuer ./build/
#
#images/enum:
#	docker build -t key-crawler-enum -f enum/Dockerfile ./enum
#
#aws/signin:
#	`aws ecr get-login --no-include-email --region us-east-1`
#
#aws/push-recorder: images/recorder aws/signin
#	docker tag key-crawler-recorder:latest 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-recorder:$(VERSION)
#	docker push 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-recorder:$(VERSION)
#
#aws/push-queuer: images/queuer aws/signin
#	docker tag key-crawler-queuer:latest 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-queuer:$(VERSION)
#	docker push 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-queuer:$(VERSION)
#
#aws/push-enum: images/enum aws/signin
#	docker tag key-crawler-enum:latest 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-enum:$(VERSION)
#	docker push 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-enum:$(VERSION)
#
#aws/push-all: aws/push-recorder aws/push-queuer aws/push-enum
#
#compose/run/enum: compose/build/enum
#	docker-compose run enum
#
#compose/build/enum:
#	docker-compose build enum
#
#release: clean images/recorder images/queuer git/tag-version aws/push-all
#
#git/tag-version:
#	git tag $(VERSION)
#	git push --tags
#
#tf/apply/dev:
#	terraform apply -var-file=infra/dev.tfvars infra/
#
#clean:
#	rm -f ./build/*
#	docker rm queuer-build recorder-build enum-enum-batch-gen-build; true
