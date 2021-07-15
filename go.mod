module github.com/loginradius/lr-cli

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.2.12
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/cli/safeexec v1.0.0
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/kr/text v0.2.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/rs/cors v1.7.0
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/spf13/cobra v1.1.3
	github.com/tklauser/go-sysconf v0.3.5 // indirect
)

replace github.com/loginradius/lr-cli => ../lr-cli
