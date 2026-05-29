package com.consumable.tracker

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.navigation.compose.rememberNavController
import com.consumable.tracker.ui.navigation.AppNavGraph
import com.consumable.tracker.ui.theme.*

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme(
                colorScheme = lightColorScheme(
                    primary = Primary,
                    onPrimary = androidx.compose.ui.graphics.Color.White,
                    background = BgPage,
                    surface = androidx.compose.ui.graphics.Color.White,
                    error = Danger,
                )
            ) {
                val navController = rememberNavController()
                AppNavGraph(navController)
            }
        }
    }
}
