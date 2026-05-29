package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import com.consumable.tracker.data.api.*
import com.consumable.tracker.data.local.Preferences
import kotlinx.coroutines.launch
import java.time.LocalDate
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(navController: NavHostController) {
    val scope = rememberCoroutineScope()
    var offices by remember { mutableStateOf<List<OfficeResponse>>(emptyList()) }
    var consumables by remember { mutableStateOf<List<ConsumableResponse>>(emptyList()) }
    var records by remember { mutableStateOf<List<RecordResponse>>(emptyList()) }
    var selectedOffice by remember { mutableStateOf<Int?>(null) }
    var selectedConsumable by remember { mutableStateOf<Int?>(null) }
    var quantity by remember { mutableStateOf(1) }
    var note by remember { mutableStateOf("") }
    var loading by remember { mutableStateOf(false) }

    LaunchedEffect(Unit) {
        try {
            val api = ApiService.create()
            offices = api.getOffices()
            consumables = api.getConsumables()
            records = api.getRecords(mapOf("page" to "1", "page_size" to "20")).data
        } catch (_: Exception) {}
    }

    LazyColumn(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        item {
            Card(modifier = Modifier.fillMaxWidth()) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("快速录入", style = MaterialTheme.typography.titleMedium)
                    Spacer(Modifier.height(12.dp))
                    if (offices.isEmpty() || consumables.isEmpty()) {
                        Text("请先在设置中添加办公室和耗材", color = MaterialTheme.colorScheme.error)
                    } else {
                        ExposedDropdownMenuBox(expanded = false, onExpandedChange = {}) {
                            OutlinedTextField(value = offices.find { it.id == selectedOffice }?.room_number ?: "选择办公室", onValueChange = {}, readOnly = true, modifier = Modifier.fillMaxWidth())
                        }
                        Spacer(Modifier.height(8.dp))
                        // Simplified - full selection would use dropdown menus
                        OutlinedTextField(value = "已选办公室: ${offices.find { it.id == selectedOffice }?.room_number ?: "无"}", onValueChange = {}, readOnly = true, modifier = Modifier.fillMaxWidth())
                        Spacer(Modifier.height(8.dp))
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text("数量: $quantity", modifier = Modifier.weight(1f))
                            IconButton(onClick = { if (quantity > 1) quantity-- }) { Text("-") }
                            IconButton(onClick = { quantity++ }) { Text("+") }
                        }
                        Spacer(Modifier.height(8.dp))
                        OutlinedTextField(value = note, onValueChange = { note = it }, label = { Text("备注") }, singleLine = true, modifier = Modifier.fillMaxWidth())
                        Spacer(Modifier.height(12.dp))
                        Button(onClick = {
                            scope.launch {
                                if (selectedOffice == null || selectedConsumable == null) return@launch
                                loading = true
                                try {
                                    val api = ApiService.create()
                                    api.createRecord(CreateRecordRequest(
                                        office_id = selectedOffice!!,
                                        consumable_id = selectedConsumable!!,
                                        quantity = quantity,
                                        usage_date = LocalDate.now().format(DateTimeFormatter.ISO_LOCAL_DATE),
                                        note = note
                                    ))
                                    quantity = 1; note = ""
                                    records = api.getRecords(mapOf("page" to "1", "page_size" to "20")).data
                                } catch (_: Exception) {} finally { loading = false }
                            }
                        }, modifier = Modifier.fillMaxWidth(), enabled = !loading) { Text("录入") }
                    }
                }
            }
            Spacer(Modifier.height(16.dp))
            Text("最近记录", style = MaterialTheme.typography.titleMedium)
            Spacer(Modifier.height(8.dp))
        }
        items(records) { record ->
            Card(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
                Row(modifier = Modifier.padding(12.dp)) {
                    Column(modifier = Modifier.weight(1f)) {
                        Text("${record.office_name} - ${record.consumable_name}", style = MaterialTheme.typography.bodyLarge)
                        Text("${record.usage_date} | ${record.quantity}个 | ${record.note}", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                    }
                }
            }
        }
    }
}
