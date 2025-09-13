import { useEffect, useState, useRef, useCallback } from "react";

interface UseWebsocketOptions {
  url?: string;
  onMessage?: (data: any) => void;
  onError?: (error: Event) => void;
  onOpen?: () => void;
  onClose?: () => void;
}

export function useWebsocket(options: UseWebsocketOptions = {}) {
  const {
    url = "ws://192.168.1.135:9090/ws",
    onMessage,
    onError,
    onOpen,
    onClose,
  } = options;

  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const messageListeners = useRef<Map<string, (data: any) => void>>(new Map());

  // Store callbacks in refs to avoid recreating the effect
  const onMessageRef = useRef(onMessage);
  const onErrorRef = useRef(onError);
  const onOpenRef = useRef(onOpen);
  const onCloseRef = useRef(onClose);

  // Update refs when callbacks change
  useEffect(() => {
    onMessageRef.current = onMessage;
    onErrorRef.current = onError;
    onOpenRef.current = onOpen;
    onCloseRef.current = onClose;
  });

  const sendMessage = useCallback((data: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const message = typeof data === "string" ? data : JSON.stringify(data);
      wsRef.current.send(message);
      return true;
    }
    return false;
  }, []);

  const waitForMessage = useCallback(
    (messageType: string, timeout: number = 5000): Promise<any> => {
      return new Promise((resolve, reject) => {
        const timeoutId = setTimeout(() => {
          messageListeners.current.delete(messageType);
          reject(new Error(`Timeout waiting for message type: ${messageType}`));
        }, timeout);

        messageListeners.current.set(messageType, (data) => {
          clearTimeout(timeoutId);
          messageListeners.current.delete(messageType);
          resolve(data);
        });
      });
    },
    []
  );

  const onMessageListener = useCallback(
    (messageType: string, callback: (data: any) => void) => {
      messageListeners.current.set(messageType, callback);

      // Return cleanup function
      return () => {
        messageListeners.current.delete(messageType);
      };
    },
    []
  );

  useEffect(() => {
    const websocket = new WebSocket(url);
    wsRef.current = websocket;

    websocket.onopen = () => {
      setIsConnected(true);
      onOpenRef.current?.();
    };

    websocket.onclose = () => {
      setIsConnected(false);
      onCloseRef.current?.();
    };

    websocket.onerror = (error) => {
      console.error("WebSocket error:", error);
      setIsConnected(false);
      onErrorRef.current?.(error);
    };

    websocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        // Call general onMessage callback
        onMessageRef.current?.(data);

        // Call specific message type listeners
        if (data.type && messageListeners.current.has(data.type)) {
          messageListeners.current.get(data.type)?.(data);
        }
      } catch (error) {
        // Handle non-JSON messages
        onMessageRef.current?.(event.data);
      }
    };

    return () => {
      if (
        wsRef.current?.readyState === WebSocket.OPEN ||
        wsRef.current?.readyState === WebSocket.CONNECTING
      ) {
        wsRef.current.close(1000, "Component unmounted");
      }
      wsRef.current = null;
      setIsConnected(false);
      messageListeners.current.clear();
    };
  }, [url]); // Only depend on url, not the callback functions

  return {
    isConnected,
    sendMessage,
    waitForMessage,
    onMessageListener,
    ws: wsRef.current,
  };
}
