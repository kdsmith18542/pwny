import { useEffect, useRef, useState, useCallback } from "react"

export function useWebSocket(url: string) {
  const [isConnected, setIsConnected] = useState(false)
  const ws = useRef<WebSocket | null>(null)
  const [lastMessage, setLastMessage] = useState<any>(null)

  useEffect(() => {
    const socket = new WebSocket(url)
    ws.current = socket

    socket.onopen = () => setIsConnected(true)
    socket.onclose = () => setIsConnected(false)
    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        setLastMessage(data)
      } catch (err) {
        setLastMessage({ type: "raw", data: event.data })
      }
    }

    return () => {
      socket.close()
    }
  }, [url])

  const sendMessage = useCallback((data: any) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(data))
    }
  }, [])

  return { isConnected, lastMessage, sendMessage }
}
