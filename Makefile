PY_FILES := $(shell git diff-index --name-only HEAD | grep .py)

reqs: 
	pip install -r requirements.txt

test:
	server/manage.py test

lint:
	pep8 .
	pyflakes .

precommit: test
	pep8 $(PY_FILES)
	pyflakes $(PY_FILES)
	echo "Working in $(shell pwd)"
	./bin/test_changes.sh $(PY_FILES)

clean: 
	find . -name '*.orig' -delete
	find . -name '*.swp' -delete
