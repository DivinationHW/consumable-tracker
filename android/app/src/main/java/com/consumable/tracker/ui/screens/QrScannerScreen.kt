package com.consumable.tracker.ui.screens

import androidx.camera.core.CameraSelector
import androidx.camera.core.ImageAnalysis
import androidx.camera.core.Preview
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.camera.view.PreviewView
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.viewinterop.AndroidView
import androidx.core.content.ContextCompat
import com.consumable.tracker.ui.theme.SidebarBg
import com.google.mlkit.vision.barcode.BarcodeScanning
import com.google.mlkit.vision.barcode.common.Barcode
import com.google.mlkit.vision.common.InputImage
import java.util.concurrent.Executors

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun QrScannerScreen(onScanned: (String) -> Unit) {
    var scanned by remember { mutableStateOf(false) }

    Scaffold(
        topBar = { TopAppBar(title = { Text("扫描二维码") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = SidebarBg, titleContentColor = androidx.compose.ui.graphics.Color.White)) }
    ) { padding ->
        Box(modifier = Modifier.fillMaxSize().padding(padding)) {
            val context = LocalContext.current
            AndroidView(
                factory = { ctx ->
                    val previewView = PreviewView(ctx)
                    val cameraProviderFuture = ProcessCameraProvider.getInstance(ctx)

                    cameraProviderFuture.addListener({
                        val cameraProvider = cameraProviderFuture.get()
                        val preview = Preview.Builder().build().also { it.surfaceProvider = previewView.surfaceProvider }
                        val cameraSelector = CameraSelector.DEFAULT_BACK_CAMERA

                        val scanner = BarcodeScanning.getClient()
                        val analysisUseCase = ImageAnalysis.Builder()
                            .setBackpressureStrategy(ImageAnalysis.STRATEGY_KEEP_ONLY_LATEST)
                            .build()
                        analysisUseCase.setAnalyzer(Executors.newSingleThreadExecutor()) { imageProxy ->
                            if (scanned) {
                                imageProxy.close()
                                return@setAnalyzer
                            }
                            @Suppress("DEPRECATION")
                            val inputImage = InputImage.fromByteBuffer(
                                imageProxy.planes[0].buffer,
                                imageProxy.width,
                                imageProxy.height,
                                imageProxy.format,
                                imageProxy.imageInfo.rotationDegrees
                            )
                            scanner.process(inputImage)
                                .addOnSuccessListener { barcodes ->
                                    for (barcode in barcodes) {
                                        barcode.rawValue?.let { value ->
                                            if (!scanned) {
                                                scanned = true
                                                onScanned(value)
                                            }
                                        }
                                    }
                                }
                                .addOnCompleteListener { imageProxy.close() }
                        }

                        cameraProvider.unbindAll()
                        cameraProvider.bindToLifecycle(
                            ctx as androidx.lifecycle.LifecycleOwner,
                            cameraSelector,
                            preview,
                            analysisUseCase
                        )
                    }, ContextCompat.getMainExecutor(ctx))

                    previewView
                },
                modifier = Modifier.fillMaxSize()
            )
        }
    }
}
