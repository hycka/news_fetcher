build:
	rm -rf dist/ \
	       	&& go build -o ./dist/yt_fetcher_server cmd/yt_fetcher/server/server.go \
		&& go build -o ./dist/yt_fetcher_manager cmd/yt_fetcher/manager/manager.go \
	       	&& go build -o ./dist/yt_fetcher_jobs cmd/yt_fetcher/jobs/jobs.go

mysql:
	docker start yt_fetcher

run:
	./dist/yt_fetcher_server

manage:
	./dist/yt_fetcher_manager

job:
	./dist/yt_fetcher_jobs

package:
	tar -czvf dist/yt_fetcher_server.latest.tar.gz dist/yt_fetcher_server \
		&& tar -czvf dist/yt_fetcher_manager.latest.tar.gz dist/yt_fetcher_manager \
		&& tar -czvf dist/yt_fetcher_jobs.latest.tar.gz dist/yt_fetcher_jobs
