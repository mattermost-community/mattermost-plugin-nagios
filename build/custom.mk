# Include custom targets and environment variables here

.PHONY: nagios
nagios:
	cd dev && docker-compose up
