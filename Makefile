.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmocks -destination=webook/internal/service/mocks/user_mock.go
	@mockgen -source=webook/internal/service/article.go -package=svcmocks -destination=webook/internal/service/mocks/article_mock.go
	@mockgen -source=webook/internal/service/code.go -package=svcmocks -destination=webook/internal/service/mocks/code_mock.go
	@mockgen -source=webook/internal/repository/user.go -package=repomocks -destination=webook/internal/repository/mocks/user_mock.go
	@mockgen -source=webook/internal/repository/code.go -package=repomocks -destination=webook/internal/repository/mocks/code_mock.go
	@mockgen -source=webook/internal/repository/article/article.go -package=artrepomocks -destination=webook/internal/repository/article/mocks/article_mock.go
	@mockgen -source=webook/internal/repository/article/article_author.go -package=artrepomocks -destination=webook/internal/repository/article/mocks/article_author_mock.go
	@mockgen -source=webook/internal/repository/article/article_reader.go -package=artrepomocks -destination=webook/internal/repository/article/mocks/article_reader_mock.go
	@mockgen -source=webook/pkg/ratelimit/types.go -package=ratelimitmocks -destination=webook/pkg/ratelimit/mocks/ratelimit_mock.go
	@go mod tidy