Nagios Core API client
======================

This is a Go client for the Nagios Core API that can communicate with the Nagios Core JSON CGIs.
It's the programmatic equivalent of the `<nagios>/nagios/jsonquery.html`.

## Installation

`go get -u github.com/ulumuri/go-nagios/nagios`

## Support

- [x] Archive JSON CGI
    - [x] alertcount
    - [x] alertlist
    - [x] notificationcount
    - [x] notificationlist
    - [ ] statechangelist
    - [ ] availability
- [ ] Object JSON CGI
    - [ ] hostcount
    - [ ] hostlist
    - [ ] host
    - [ ] hostgroupcount
    - [ ] hostgrouplist
    - [ ] hostgroup
    - [ ] servicecount
    - [ ] servicelist
    - [ ] service
    - [ ] servicegroupcount
    - [ ] servicegrouplist
    - [ ] servicegroup
    - [ ] contactcount
    - [ ] contactlist
    - [ ] contact
    - [ ] contactgroupcount
    - [ ] contactgrouplist
    - [ ] contactgroup
    - [ ] timeperiodcount
    - [ ] timeperiodlist
    - [ ] timeperiod
    - [ ] commandcount
    - [ ] commandlist
    - [ ] command
    - [ ] servicedependencycount
    - [ ] servicedependencylist
    - [ ] serviceescalationcount
    - [ ] serviceescalationlist
    - [ ] hostdependencycount
    - [ ] hostdependencylist
    - [ ] hostescalationcount
    - [ ] hostescalationlist
- [ ] Status JSON CGI
    - [x] hostcount
    - [x] hostlist
    - [x] host
    - [x] servicecount
    - [x] servicelist
    - [x] service
    - [ ] commentcount
    - [ ] commentlist
    - [ ] comment
    - [ ] downtimecount
    - [ ] downtimelist
    - [ ] downtime
    - [ ] programstatus
    - [x] performancedata

## Compatibility

Tested with Nagios Core 4.4.x
