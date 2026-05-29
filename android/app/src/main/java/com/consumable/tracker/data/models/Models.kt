package com.consumable.tracker.data.models

data class LoginRequest(val username: String, val password: String)
data class LoginResponse(val token: String, val role: String)

data class User(val id: Int, val username: String, val role: String, val created_at: String)
data class UserCreate(val username: String, val password: String, val role: String = "readonly")
data class UserUpdate(val username: String? = null, val password: String? = null, val role: String? = null)

data class Consumable(val id: Int, val name: String, val category: String = "", val unit: String = "", val stock: Int = 0, val threshold: Int = 0)
data class Office(val id: Int, val room_number: String, val department: String = "", val device_type: String = "", val device_model: String = "")
data class OfficeCreate(val room_number: String, val department: String = "", val device_type: String = "", val device_model: String = "")

data class UsageRecord(
    val id: Int, val office_id: Int, val consumable_id: Int, val quantity: Int,
    val usage_date: String, val note: String = "", val room_number: String = "",
    val consumable_name: String = "", val unit: String = ""
)
data class RecordCreate(val office_id: Int, val consumable_id: Int, val quantity: Int, val usage_date: String, val note: String = "")

data class StatsResponse(
    val office_count: Int = 0, val consumable_count: Int = 0,
    val monthly_usage: Int = 0, val pending_tickets: Int = 0
)

data class Ticket(
    val id: String, val office_id: Int, val room_number: String = "",
    val device_type: String = "", val problem_type: String = "",
    val description: String = "", val contact: String = "", val status: String = "pending",
    val created_at: String = ""
)
data class TicketCreate(val office_id: Int, val problem_type: String, val description: String = "", val contact: String = "")
data class TicketComplete(val consumable_used: String = "", val consumable_quantity: Int = 0, val handle_note: String = "")

data class ProblemType(val id: Int, val name: String, val device_type: String = "", val is_default: Boolean = false, val sort_order: Int = 0)
data class QRCode(val id: Int, val code: String, val office_id: Int? = null, val device_type: String = "", val device_model: String = "", val created_at: String = "")
data class BackupFile(val name: String, val size: Long, val date: String)
data class Note(val id: Int, val office_id: Int, val content: String, val created_at: String = "", val room_number: String = "")

data class DeviceModelsResponse(val models: List<String>)
data class ErrorResponse(val error: String)
