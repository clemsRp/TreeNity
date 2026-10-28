```Makefile
UTILS_DIR := ./utils

branch:
ifndef NAME
	@echo "'NAME' parameter is missing"
else
	@RESULT=$$(sh $(UTILS_DIR)/check_branch_name.sh "$(NAME)"); \
	if [ -z "$$RESULT" ]; then \
		git checkout -b $(NAME); \
		echo "New branch created successfully"; \
	else \
		echo "Validation error : $$RESULT"; \
		exit 1; \
	fi
endif

commit:
ifndef NAME
	@echo "'NAME' parameter is missing"
else
	@RESULT=$$(sh $(UTILS_DIR)/check_commit_name.sh "$(NAME)"); \
	if [ -z "$$RESULT" ]; then \
		git commit -m "$(NAME)"; \
		echo "Commit successfully"; \
	else \
		echo "Validation error : $$RESULT"; \
		exit 1; \
	fi
endif

sync:
	@git fetch origin main; \
	BEHIND=$$(git rev-list --count HEAD..origin/main); \
	if [ "$$BEHIND" -gt 0 ]; then \
		echo "Your branch is $$BEHIND commit(s) late"; \
		git rebase origin/main; \
		echo "Rebased main on your branch"; \
	fi; \
	echo "Branch is up to date"

push: sync
	@BRANCH=$$(git branch --show-current); \
	git push origin $$BRANCH; \
	echo "Push $$BRANCH successfully"

.PHONY: branch commit sync push
```