migrate:
	@echo "migrating all table...."
	@go run ./migration/models.go ./migration/mock_data.go ./migration/migrate.go

git-cycle:
	@echo "git cycle...."
	@git add .
	@git commit -m "update new feature..."
	@git pull origin main
	@git push origin main