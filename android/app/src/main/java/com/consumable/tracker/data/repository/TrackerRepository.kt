package com.consumable.tracker.data.repository

import com.consumable.tracker.data.api.ApiClient
import com.consumable.tracker.data.models.*

class TrackerRepository {
    private val api get() = ApiClient.service

    suspend fun login(username: String, password: String): LoginResponse {
        val res = api.login(LoginRequest(username, password))
        ApiClient.setToken(res.token)
        return res
    }

    suspend fun getConsumables() = api.getConsumables()
    suspend fun getOffices() = api.getOffices()
    suspend fun createOffice(body: OfficeCreate) = api.createOffice(body)
    suspend fun updateOffice(id: Int, body: Map<String, String>) = api.updateOffice(id, body)
    suspend fun deleteOffice(id: Int) = api.deleteOffice(id)
    suspend fun getRecords(params: Map<String, String> = emptyMap()) = api.getRecords(params)
    suspend fun createRecord(body: RecordCreate) = api.createRecord(body)
    suspend fun deleteRecord(id: Int) = api.deleteRecord(id)
    suspend fun getStats() = api.getStats()
    suspend fun getTickets(status: String? = null) = api.getTickets(status)
    suspend fun updateTicketStatus(id: String, status: String) = api.updateTicketStatus(id, mapOf("status" to status))
    suspend fun completeTicket(id: String, body: TicketComplete) = api.completeTicket(id, body)
    suspend fun deleteTicket(id: String) = api.deleteTicket(id)
    suspend fun submitTicket(body: TicketCreate) = api.submitTicket(body)
    suspend fun getTicketStatus(id: String) = api.getTicketStatus(id)
    suspend fun getQRCodes() = api.getQRCodes()
    suspend fun getDeviceModels() = api.getDeviceModels()
    suspend fun getBackupList() = api.getBackupList()
    suspend fun createBackup() = api.createBackup()
    suspend fun deleteBackup(name: String) = api.deleteBackup(name)
}
