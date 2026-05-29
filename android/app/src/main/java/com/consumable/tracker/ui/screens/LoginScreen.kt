package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import com.consumable.tracker.data.api.ApiService
import com.consumable.tracker.data.api.LoginRequest
import com.consumable.tracker.data.local.Preferences
import com.consumable.tracker.ui.navigation.Screen
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LoginScreen(navController: NavHostController) {
    var username by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf("") }
    val scope = rememberCoroutineScope()

    Column(
        modifier = Modifier.fillMaxSize().padding(32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text("耗材管理系统", style = MaterialTheme.typography.headlineLarge)
        Spacer(Modifier.height(32.dp))
        OutlinedTextField(value = username, onValueChange = { username = it; error = "" }, label = { Text("用户名") }, singleLine = true, modifier = Modifier.fillMaxWidth())
        Spacer(Modifier.height(16.dp))
        OutlinedTextField(value = password, onValueChange = { password = it; error = "" }, label = { Text("密码") }, singleLine = true, visualTransformation = PasswordVisualTransformation(), modifier = Modifier.fillMaxWidth())
        Spacer(Modifier.height(8.dp))
        if (error.isNotEmpty()) { Text(error, color = MaterialTheme.colorScheme.error); Spacer(Modifier.height(8.dp)) }
        Button(onClick = {
            scope.launch {
                loading = true; error = ""
                try {
                    val api = ApiService.create()
                    val res = api.login(LoginRequest(username, password))
                    Preferences.saveLogin(res.token, res.user_id, res.username, res.role)
                    navController.navigate(Screen.Home.route) { popUpTo(0) }
                } catch (e: Exception) { error = "登录失败: ${e.message}" } finally { loading = false }
            }
        }, modifier = Modifier.fillMaxWidth().height(48.dp), enabled = !loading) {
            if (loading) CircularProgressIndicator(modifier = Modifier.size(24.dp), color = MaterialTheme.colorScheme.onPrimary) else Text("登录")
        }
    }
}
