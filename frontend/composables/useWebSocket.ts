// WebSocket composable for real-time messaging
import { ref, onUnmounted } from 'vue'

interface WebSocketMessage {
    type: string // new_message, typing, read, presence
    conversation_id?: string
    sender_id?: string
    content?: any
}

const socket = ref<WebSocket | null>(null)
const connected = ref(false)
const messages = ref<WebSocketMessage[]>([])
const onMessageCallbacks: ((msg: WebSocketMessage) => void)[] = []

export function useWebSocket() {
    const getWsUrl = () => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const apiHost = 'api.app.tech-afm.com'
        return `${protocol}//${apiHost}/messaging-service/ws/chat`
    }

    const connect = (userId: string) => {
        if (socket.value?.readyState === WebSocket.OPEN) {
            console.log('WebSocket already connected')
            return
        }

        const url = `${getWsUrl()}?user_id=${userId}`
        console.log('WebSocket connecting to:', url)

        socket.value = new WebSocket(url)

        socket.value.onopen = () => {
            console.log('WebSocket connected')
            connected.value = true
        }

        socket.value.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data) as WebSocketMessage
                messages.value.push(data)

                // Call all registered callbacks
                onMessageCallbacks.forEach(cb => cb(data))
            } catch (e) {
                console.error('WebSocket message parse error:', e)
            }
        }

        socket.value.onclose = (event) => {
            console.log('WebSocket disconnected:', event.code, event.reason)
            connected.value = false

            // Auto-reconnect after 3 seconds if not intentionally closed
            if (event.code !== 1000) {
                setTimeout(() => {
                    console.log('WebSocket reconnecting...')
                    connect(userId)
                }, 3000)
            }
        }

        socket.value.onerror = (error) => {
            console.error('WebSocket error:', error)
        }
    }

    const disconnect = () => {
        if (socket.value) {
            socket.value.close(1000, 'User disconnected')
            socket.value = null
            connected.value = false
        }
    }

    const send = (message: WebSocketMessage) => {
        if (socket.value?.readyState === WebSocket.OPEN) {
            socket.value.send(JSON.stringify(message))
        } else {
            console.warn('WebSocket not connected, cannot send message')
        }
    }

    const sendTyping = (conversationId: string) => {
        send({
            type: 'typing',
            conversation_id: conversationId
        })
    }

    const sendRead = (conversationId: string, messageId: string) => {
        send({
            type: 'read',
            conversation_id: conversationId,
            content: { message_id: messageId }
        })
    }

    const onMessage = (callback: (msg: WebSocketMessage) => void) => {
        onMessageCallbacks.push(callback)
    }

    const offMessage = (callback: (msg: WebSocketMessage) => void) => {
        const index = onMessageCallbacks.indexOf(callback)
        if (index > -1) {
            onMessageCallbacks.splice(index, 1)
        }
    }

    // Cleanup on unmount
    onUnmounted(() => {
        // Don't disconnect on unmount - keep connection alive
        // disconnect()
    })

    return {
        socket,
        connected,
        messages,
        connect,
        disconnect,
        send,
        sendTyping,
        sendRead,
        onMessage,
        offMessage
    }
}

// Singleton instance for global WebSocket connection
let globalWs: ReturnType<typeof useWebSocket> | null = null

export function useGlobalWebSocket() {
    if (!globalWs) {
        globalWs = useWebSocket()
    }
    return globalWs
}
