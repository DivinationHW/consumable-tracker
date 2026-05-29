package com.consumable.tracker.data.api

import com.consumable.tracker.data.models.*
import retrofit2.http.*

interface ApiService {
    @POST("api/auth/login")
    suspend fun login(@Body request: LoginRequest): LoginResponse

    @PUT("api/auth/password")
    suspend fun changePassword(@Body body: Map<String, String>)

    @GET("api/users")
    suspend fun getUsers(): List<User>

    @POST("api/users")
    suspend fun createUser(@Body body: UserCreate)

    @PUT("api/users/{id}")
    suspend fun updateUser(@Path("id") id: Int, @Body body: UserUpdate)

    @DELETE("api/users/{id}")
    suspend fun deleteUser(@Path("id") id: Int)

    @GET("api/consumables")
    suspend fun getConsumables(): List<Consumable>

    @GET("api/offices")
    suspend fun getOffices(): List<Office>

    @POST("api/offices")
    suspend fun createOffice(@Body body: OfficeCreate)

    @PUT("api/offices/{id}")
    suspend fun updateOffice(@Path("id") id: Int, @Body body: Map<String, String>)

    @DELETE("api/offices/{id}")
    suspend fun deleteOffice(@Path("id") id: Int)

    @GET("api/records")
    suspend fun getRecords(@QueryMap params: Map<String, String> = emptyMap()): List<UsageRecord>

    @POST("api/records")
    suspend fun createRecord(@Body body: RecordCreate)

    @DELETE("api/records/{id}")
    suspend fun deleteRecord(@Path("id") id: Int)

    @GET("api/stats/summary")
    suspend fun getStats(): StatsResponse

    @GET("api/tickets")
    suspend fun getTickets(@Query("status") status: String? = null): List<Ticket>

    @POST("api/tickets/{id}/status")
    suspend fun updateTicketStatus(@Path("id") id: String, @Body body: Map<String, String>)

    @POST("api/tickets/{id}/complete")
    suspend fun completeTicket(@Path("id") id: String, @Body body: TicketComplete)

    @DELETE("api/tickets/{id}")
    suspend fun deleteTicket(@Path("id") id: String)

    @POST("ticket")
    suspend fun submitTicket(@Body body: TicketCreate): Ticket

    @GET("ticket/{id}")
    suspend fun getTicketStatus(@Path("id") id: String): Ticket

    @GET("api/qrcodes")
    suspend fun getQRCodes(): List<QRCode>

    @GET("api/device-models")
    suspend fun getDeviceModels(): DeviceModelsResponse

    @GET("api/backup/list")
    suspend fun getBackupList(): List<BackupFile>

    @POST("api/backup/now")
    suspend fun createBackup()

    @DELETE("api/backup/{name}")
    suspend fun deleteBackup(@Path("name") name: String)
}
