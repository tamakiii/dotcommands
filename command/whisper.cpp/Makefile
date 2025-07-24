.PHONY: help depend independ

help:
	@cat $(firstword $(MAKEFILE_LIST))

depend:
	make -f dependency.mk setup
	make -f dependency.mk build

independ:
	make -f dependency.mk clean
	make -f dependency.mk teardown
