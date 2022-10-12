all: build

build:
	@cd web && yarn build
	@rsync -e ssh -av --delete --exclude=.idea,.DS_Store server/ wechaty:/docker/view
	@ssh wechaty "cd /docker/view ; xgo --targets=linux/amd64 /docker/view ; zip -r view.zip view-linux-amd64 README.md view.db static"
	@scp wechaty:/docker/view/view.zip .

init:
	@cd web && yarn install
	@cd server && go mod tidy
