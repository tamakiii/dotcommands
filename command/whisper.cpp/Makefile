.PHONY: help depend independ setup teardown

help:
	@cat $(firstword $(MAKEFILE_LIST))

depend:
	make -f dependency.mk setup
	make -f dependency.mk build

independ:
	make -f dependency.mk clean
	make -f dependency.mk teardown

setup: \
	dist \
	dist/mp3

teardown:
	rm -rf dist/mp3
	rm -rf dist

dist:
	-mkdir $@

dist/mp3:
	-mkdir $@
