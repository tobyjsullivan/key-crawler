all: recorder queuer

recorder: build/recorder
	docker build -t key-crawler-recorder -f ./recorder/Dockerfile .

build/recorder:
	docker build -t recorder-build -f ./recorder/Dockerfile.build .
	docker create --rm --name recorder-build recorder-build
	docker cp recorder-build:/go/bin/recorder ./build/

queuer: build/queuer
	docker build -t key-crawler-queuer -f ./queuer/Dockerfile .

build/queuer:
	docker build -t queuer-build -f ./queuer/Dockerfile.build .
	docker create --rm --name queuer-build queuer-build
	docker cp queuer-build:/go/bin/queuer ./build/

aws-signin:
	`aws ecr get-login --no-include-email --region us-east-1`

push-recorder: recorder aws-signin
	docker tag key-crawler-recorder:latest 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-recorder:latest
	docker push 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-recorder:latest

push-queuer: queuer aws-signin 
	docker tag key-crawler-queuer:latest 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-queuer:latest
	docker push 110303772622.dkr.ecr.us-east-1.amazonaws.com/key-crawler-queuer:latest

push-all: push-recorder push-queuer

clean:
	rm -f ./build/*
	docker rm queuer-build recorder-build; true
