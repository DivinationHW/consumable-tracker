package com.consumable.tracker.ui.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import com.consumable.tracker.ui.screens.*

sealed class Screen(val route: String) {
    object Login : Screen("login")
    object Home : Screen("home")
    object Records : Screen("records")
    object Stats : Screen("stats")
    object Tickets : Screen("tickets")
    object QRCodes : Screen("qrcodes")
    object Notes : Screen("notes")
    object Settings : Screen("settings")
    object Backup : Screen("backup")
}

@Composable
fun NavGraph(navController: NavHostController, isAdmin: Boolean) {
    NavHost(navController, startDestination = Screen.Login.route) {
        composable(Screen.Login.route) { LoginScreen(navController) }
        composable(Screen.Home.route) { HomeScreen(navController) }
        composable(Screen.Records.route) { RecordsScreen(navController) }
        composable(Screen.Stats.route) { StatsScreen() }
        composable(Screen.Tickets.route) { TicketListScreen() }
        composable(Screen.QRCodes.route) { QRCodeScreen(navController) }
        composable(Screen.Notes.route) { NotesScreen() }
        composable(Screen.Settings.route) { SettingsScreen(navController) }
        composable(Screen.Backup.route) { BackupsScreen() }
    }
}
