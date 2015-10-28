# monlite
Monitor services and send e-mail warning if something is wrong.
Supported url schemes are: mongodb, couch, http, imap, ldap, mysql, smtp, dns, tcp, udp and unix.

# configuration
The format of the configuration file is ini. See example:

# monlite.ini
```
[log]
file=monlite.log
level=debug

[mail]
smtp=smtp.gmail.com:587
account=name@gmail.com
password=secret
from=name@gmail.com
to=foo@gmail.com
helo=
timeout=60

[service]
timeout=300
periode=300
fails=2
sleep=3600

[service.http]
url=https://www.google.com

[service.smtp]
url=smtp://smtp.gmail.com:25

[service.imap]
url=imap://imap.gmail.com:993
```
