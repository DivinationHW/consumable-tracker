sealed class Screen(val route: String, val title: String) {
    data object Login : Screen("login", "登录")
    data object Dashboard : Screen("dashboard", "仪表盘")
    data object Records : Screen("records", "使用记录")
    data object Stats : Screen("stats", "数据统计")
    data object Tickets : Screen("tickets", "报修工单")
    data object TicketForm : Screen("ticket_form/{code}", "提交报修") {
        fun createRoute(code: String) = "ticket_form/$code"
    }
    data object QrScanner : Screen("qr_scanner", "扫描二维码")
    data object Settings : Screen("settings", "设置")
}
