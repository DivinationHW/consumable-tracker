import { ref, onUnmounted } from 'vue'

export function useWebSocket() {
  const connected = ref(false)
  let ws: WebSocket | null = null
  const handlers = new Map<string, (data: any) => void>()

  function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = `${protocol}//${window.location.host}/ws`
    const token = localStorage.getItem('token')
    ws = new WebSocket(url)
    ws.onopen = () => { connected.value = true }
    ws.onclose = () => { connected.value = false }
    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        const handler = handlers.get(msg.type)
        if (handler) handler(msg)
      } catch {}
    }
  }

  function on(type: string, handler: (data: any) => void) {
    handlers.set(type, handler)
  }

  function disconnect() {
    ws?.close()
    ws = null
    connected.value = false
  }

  onUnmounted(disconnect)

  return { connect, disconnect, on, connected }
}
