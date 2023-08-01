start powershell.exe -NoExit -Command "cd src\8080_UserService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8081_SiteManagerService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8082_DeviceManagerService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8083_DashboardService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8084_ReportService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8085_AlertService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8086_LocalLiveService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8087_RuleService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8088_AnalyticsService; go run main.go"
start powershell.exe -NoExit -Command "cd src\8100_StaticService; go run main.go"
cd src\nginx-server
nginx.exe