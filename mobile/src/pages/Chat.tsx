import { useNavigation } from "@react-navigation/native";
import { Button, StyleSheet, Text, TextInput, View } from "react-native";
import { useWebsocket } from "../hooks/useWebsocket";
import { useState } from "react";

export function Chat() {
  const [message, setMessage] = useState<string>("");
  const [serverMessage, setServerMessage] = useState<string | null>(null);
  const [count, setCount] = useState(0);
  const { isConnected, sendMessage } = useWebsocket({
    onOpen: () => {
      console.log("WebSocket connected");
      sendMessage({ type: "greet", payload: "Hello, server!" });
    },
    onError: (error) => {
      console.error("WebSocket error:", error);
    },
    onMessage: (data) => {
      console.log("Received message from server:", data);
      setServerMessage(typeof data === "string" ? data : JSON.stringify(data));
      setCount((prev) => prev + 1);
    },
  });
  const navigation = useNavigation();

  function handleSendMessage() {
    sendMessage({ type: "chat", payload: message });
    setMessage("");
  }

  return (
    <View style={styles.container}>
      <Text>This is the Chat screen</Text>
      <Text>Status: {isConnected ? "Connected" : "Not connected"}</Text>
      <Text>Messages received: {count}</Text>
      {serverMessage && <Text>Last message: {serverMessage}</Text>}

      <TextInput
        placeholder="Type your message"
        value={message}
        onChangeText={setMessage}
      />
      <Button title="Send" onPress={handleSendMessage} />

      <Button
        title="Back to home"
        onPress={() => navigation.navigate("Home" as never)}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#fff",
    alignItems: "center",
    justifyContent: "center",
  },
});
