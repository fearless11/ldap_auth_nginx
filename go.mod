module ldap_auth_nginx

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => github.com/golang/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 => github.com/golang/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/mod v0.1.1-0.20191105210325-c90efee705ee => github.com/golang/mod v0.1.1-0.20191105210325-c90efee705ee
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 => github.com/golang/net v0.0.0-20190404232315-eb5bcb51f2a3
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 => github.com/golang/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => github.com/golang/sys v0.0.0-20190215142949-d0b11bdaac8a
	golang.org/x/sys v0.0.0-20190412213103-97732733099d => github.com/golang/sys v0.0.0-20190412213103-97732733099d
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/tools v0.0.0-20200117065230-39095c1d176c => github.com/golang/tools v0.0.0-20200117065230-39095c1d176c
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 => github.com/golang/xerrors v0.0.0-20191011141410-1b5146add898
	gopkg.in/asn1-ber.v1 v1.4.1 => github.com/go-asn1-ber/asn1-ber v1.4.1
)

require (
	github.com/astaxie/beego v1.12.1
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	gopkg.in/asn1-ber.v1 v1.4.1 // indirect
	gopkg.in/ldap.v2 v2.5.1
)
