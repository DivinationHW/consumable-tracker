package com.consumable.tracker

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.compose.rememberNavController
import com.consumable.tracker.data.local.Preferences
import com.consumable.tracker.ui.navigation.NavGraph
import com.consumable.tracker.ui.navigation.Screen
import com.consumable.tracker.ui.theme.TrackerTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Preferences.init(this)

        setContent {
            TrackerTheme {
                val navController = rememberNavController()
                val isLoggedIn = remember { Preferences.token.isNotEmpty() }
                val isAdmin = remember { Preferences.role == "admin" }
                val showDrawer = remember { mutableStateOf(false) }
                val currentRoute = remember { mutableStateOf("home") }

                Scaffold(
                    topBar = {
                        if (isLoggedIn) {
                            NavigationDrawerItem(
                                icon = { Icon(Icons.Default.Menu, "Menu") },
                                label = { Text(Preferences.username) },
                                selected = false,
                                onClick = { showDrawer.value = true }
                            )
                            TopAppBar(
                                title = { Text("耗材管理") },
                                navigationIcon = {
                                    IconButton(onClick = { showDrawer.value = true }) {
                                        Icon(Icons.Default.Menu, "菜单")
                                    }
                                },
                                actions = {
                                    IconButton(onClick = {
                                        Preferences.clear()
                                        navController.navigate(Screen.Login.route) { popUpTo(0) }
                                    }) {
                                        Icon(Icons.Default.ExitToApp, "退出")
                                    }
                                }
                            )
                        }
                    }
                ) { padding ->
                    if (isLoggedIn) {
                        NavigationDrawer(
                            drawerState = DrawerState(if (showDrawer.value) DrawerValue.Open else DrawerValue.Closed),
                            gesturesEnabled = showDrawer.value,
                            drawerContent = {
                                ModalDrawerSheet {
                                    Text("耗材管理系统", style = MaterialTheme.typography.headlineSmall, modifier = Modifier.padding(16.dp))
                                    Divider()
                                    NavItem("首页", Icons.Default.Home, Screen.Home.route, currentRoute, navController) { showDrawer.value = false }
                                    NavItem("记录查询", Icons.Default.List, Screen.Records.route, currentRoute, navController) { showDrawer.value = false }
                                    NavItem("统计报表", Icons.Default.BarChart, Screen.Stats.route, currentRoute, navController) { showDrawer.value = false }
                                    NavItem("工单管理", Icons.Default.Build, Screen.Tickets.route, currentRoute, navController) { showDrawer.value = false }
                                    NavItem("二维码", Icons.Default.QrCode, Screen.QRCodes.route, currentRoute, navController) { showDrawer.value = false }
                                    if (isAdmin) {
                                        NavItem("便签", Icons.Default.Note, Screen.Notes.route, currentRoute, navController) { showDrawer.value = false }
                                    }
                                    NavItem("设置", Icons.Default.Settings, Screen.Settings.route, currentRoute, navController) { showDrawer.value = false }
                                    if (isAdmin) {
                                        NavItem("备份管理", Icons.Default.Folder, Screen.Backup.route, currentRoute, navController) { showDrawer.value = false }
                                    }
                                }
                            }
                        ) {
                            Box(modifier = Modifier.padding(padding)) {
                                NavGraph(navController, isAdmin)
                            }
                        }
                    } else {
                        NavGraph(navController, isAdmin)
                    }
                }
            }
        }
    }
}

@Composable
private fun NavItem(label: String, icon: androidx.compose.ui.graphics.vector.ImageVector, route: String, currentRoute: MutableState<String>, navController: androidx.navigation.NavHostController, onClose: () -> Unit) {
    NavigationDrawerItem(
        icon = { Icon(icon, label) },
        label = { Text(label) },
        selected = currentRoute.value == route,
        onClick = {
            currentRoute.value = route
            navController.navigate(route) { popUpTo(Screen.Home.route) { saveState = true }; launchSingleTop = true; restoreState = true }
            onClose()
        },
        modifier = Modifier.padding(NavigationDrawerItemDefaults.ItemPadding)
    )
}
