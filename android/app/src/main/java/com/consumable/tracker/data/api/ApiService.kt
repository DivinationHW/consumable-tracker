package com.consumable.tracker.data.api

import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.*
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import java.util.concurrent.TimeUnit

interface ApiService {
    @POST("api/login")
    suspend fun login(@Body req: LoginRequest): LoginResponse

    @GET("api/me")
    suspend fun getMe(): UserResponse

    @GET("api/offices")
    suspend fun getOffices(): List<OfficeResponse>

    @GET("api/consumables")
    suspend fun getConsumables(): List<ConsumableResponse>

    @GET("api/records")
    suspend fun getRecords(@QueryMap params: Map<String, String>): RecordListResponse

    @POST("api/records")
    suspend fun createRecord(@Body req: CreateRecordRequest): IdResponse

    @PUT("api/records/{id}")
    suspend fun updateRecord(@Path("id") id: Int, @Body req: CreateRecordRequest): MessageResponse

    @DELETE("api/records/{id}")
    suspend fun deleteRecord(@Path("id") id: Int): MessageResponse

    @GET("api/stats")
    suspend fun getStats(): StatsResponse

    @GET("api/tickets")
    suspend fun getTickets(@Query("status") status: String? = null): List<TicketResponse>

    @POST("api/tickets/{id}/process")
    suspend fun processTicket(@Path("id") id: String, @Body req: ProcessTicketRequest): MessageResponse

    @GET("api/qrcodes")
    suspend fun getQRCodes(): List<QRCodeResponse>

    @POST("api/qrcodes")
    suspend fun createQRCode(@Body req: CreateQRRequest): IdResponse

    @DELETE("api/qrcodes/{id}")
    suspend fun deleteQRCode(@Path("id") id: Int): MessageResponse

    @GET("api/notes")
    suspend fun getNotes(): List<NoteResponse>

    @POST("api/notes")
    suspend fun createNote(@Body req: NoteRequest): IdResponse

    @PUT("api/notes/{id}")
    suspend fun updateNote(@Path("id") id: Int, @Body req: NoteRequest): MessageResponse

    @DELETE("api/notes/{id}")
    suspend fun deleteNote(@Path("id") id: Int): MessageResponse

    @GET("api/users")
    suspend fun getUsers(): List<UserListResponse>

    @POST("api/users")
    suspend fun createUser(@Body req: CreateUserRequest): IdResponse

    @PUT("api/users/{id}")
    suspend fun updateUser(@Path("id") id: Int, @Body req: CreateUserRequest): MessageResponse

    @DELETE("api/users/{id}")
    suspend fun deleteUser(@Path("id") id: Int): MessageResponse

    @POST("api/change-password")
    suspend fun changePassword(@Body req: ChangePasswordRequest): MessageResponse

    @GET("api/backups")
    suspend fun getBackups(): List<BackupResponse>

    @POST("api/backups")
    suspend fun createBackup(): MessageResponse

    @POST("api/backups/{filename}/restore")
    suspend fun restoreBackup(@Path("filename") filename: String): MessageResponse

    @DELETE("api/backups/{filename}")
    suspend fun deleteBackup(@Path("filename") filename: String): MessageResponse

    @GET("api/problem-types")
    suspend fun getProblemTypes(@Query("device_type") deviceType: String? = null): List<ProblemTypeResponse>

    @POST("api/problem-types")
    suspend fun createProblemType(@Body req: ProblemTypeRequest): IdResponse

    @PUT("api/problem-types/{id}")
    suspend fun updateProblemType(@Path("id") id: Int, @Body req: ProblemTypeRequest): MessageResponse

    @DELETE("api/problem-types/{id}")
    suspend fun deleteProblemType(@Path("id") id: Int): MessageResponse

    @GET("api/device-models")
    suspend fun getDeviceModels(): List<String>

    @GET("api/device-types")
    suspend fun getDeviceTypes(): List<String>

    companion object {
        var baseUrl = "http://10.0.2.2:8443/"

        fun create(): ApiService {
            val logging = HttpLoggingInterceptor().apply { level = HttpLoggingInterceptor.Level.BODY }
            val client = OkHttpClient.Builder()
                .addInterceptor { chain ->
                    val req = chain.request().newBuilder()
                    val token = com.consumable.tracker.data.local.Preferences.token
                    if (token.isNotEmpty()) req.addHeader("Authorization", "Bearer $token")
                    chain.proceed(req.build())
                }
                .connectTimeout(10, TimeUnit.SECONDS)
                .readTimeout(30, TimeUnit.SECONDS)
                .build()
            return Retrofit.Builder()
                .baseUrl(baseUrl)
                .client(client)
                .addConverterFactory(GsonConverterFactory.create())
                .build()
                .create(ApiService::class.java)
        }
    }
}

// Request/Response models
data class LoginRequest(val username: String, val password: String)
data class LoginResponse(val token: String, val user_id: Int, val username: String, val role: String)
data class UserResponse(val user_id: Int, val username: String, val role: String)
data class OfficeResponse(val id: Int, val room_number: String, val device_type: String, val device_model: String, val created_at: String)
data class ConsumableResponse(val id: Int, val name: String, val unit: String, val is_default: Boolean, val created_at: String)
data class RecordListResponse(val data: List<RecordResponse>, val total: Int, val page: Int)
data class RecordResponse(val id: Int, val office_id: Int, val consumable_id: Int, val quantity: Int, val usage_date: String, val note: String, val created_at: String, val updated_at: String, val office_name: String, val consumable_name: String, val device_type: String, val device_model: String)
data class CreateRecordRequest(val office_id: Int, val consumable_id: Int, val quantity: Int = 1, val usage_date: String, val note: String = "")
data class IdResponse(val id: Int)
data class MessageResponse(val message: String)
data class StatsResponse(val total_usage: Int, val current_month: Int, val top_office: String, val top_consumable: String, val by_office: List<StatsItem>, val monthly_trend: List<StatsItem>)
data class StatsItem(val label: String, val value: Int)
data class TicketResponse(val id: String, val office_id: Int, val device_type: String, val device_model: String, val problem_type: String, val description: String, val contact: String, val status: String, val consumable_used: String, val consumable_quantity: Int, val handled_by_user_id: Int?, val handle_note: String, val created_at: String, val updated_at: String, val office_name: String, val handled_by_user: String?)
data class ProcessTicketRequest(val status: String, val consumable_used: String = "", val consumable_quantity: Int = 0, val handle_note: String = "")
data class QRCodeResponse(val id: Int, val code: String, val office_id: Int?, val device_type: String, val device_model: String, val created_at: String, val office_name: String)
data class CreateQRRequest(val code: String = "", val office_id: Int? = null, val device_type: String = "", val device_model: String = "")
data class NoteResponse(val id: Int, val title: String, val content: String, val created_at: String, val updated_at: String)
data class NoteRequest(val title: String, val content: String)
data class UserListResponse(val id: Int, val username: String, val role: String, val created_at: String, val updated_at: String)
data class CreateUserRequest(val username: String, val password: String, val role: String = "readonly")
data class ChangePasswordRequest(val old_password: String, val new_password: String)
data class BackupResponse(val filename: String, val size: Long, val created_at: String)
data class ProblemTypeResponse(val id: Int, val device_type: String, val name: String, val sort_order: Int, val is_default: Boolean, val created_at: String)
data class ProblemTypeRequest(val device_type: String, val name: String, val sort_order: Int = 0)
