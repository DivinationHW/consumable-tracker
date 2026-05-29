package com.consumable.tracker.ui.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.navArgument
import com.consumable.tracker.ui.screens.*

@Composable
fun AppNavGraph(navController: NavHostController) {
    NavHost(navController = navController, startDestination = Screen.Login.route) {
        composable(Screen.Login.route) {
            LoginScreen(onLoginSuccess = {
                navController.navigate(Screen.Dashboard.route) {
                    popUpTo(Screen.Login.route) { inclusive = true }
                }
            })
        }
        composable(Screen.Dashboard.route) {
            DashboardScreen(onNavigate = { navController.navigate(it.route) })
        }
        composable(Screen.Records.route) {
            RecordsScreen()
        }
        composable(Screen.Stats.route) {
            StatsScreen()
        }
        composable(Screen.Tickets.route) {
            TicketsScreen()
        }
        composable(
            Screen.TicketForm.route,
            arguments = listOf(navArgument("code") { type = NavType.StringType })
        ) { backStackEntry ->
            val code = backStackEntry.arguments?.getString("code") ?: ""
            TicketFormScreen(code, onSubmitted = { navController.popBackStack() })
        }
        composable(Screen.QrScanner.route) {
            QrScannerScreen(onScanned = { code ->
                navController.navigate(Screen.TicketForm.createRoute(code))
            })
        }
        composable(Screen.Settings.route) {
            SettingsScreen(onLogout = {
                navController.navigate(Screen.Login.route) {
                    popUpTo(0) { inclusive = true }
                }
            })
        }
    }
}
