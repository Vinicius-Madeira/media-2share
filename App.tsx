import { StatusBar } from "expo-status-bar";
import { useState } from "react";
import * as ImagePicker from "expo-image-picker";
import { Button, Image, StyleSheet, Text, View } from "react-native";

export default function App() {
  const [image, setImage] = useState<string | null>(null);

  async function pickImage() {
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ["images"],
      allowsEditing: false,
      aspect: [4, 3],
      quality: 1,
    });
    console.log("Image result:", result);

    if (result.canceled) {
      setImage(null);
      return;
    }
    setImage(result.assets[0].uri);
  }

  return (
    <View style={styles.container}>
      <Text>Open up App.tsx to start working on your app!</Text>
      <Text>Hello World</Text>
      <StatusBar style="auto" />
      <Button title="Pick an image from your galery" onPress={pickImage} />
      {image && <Image source={{ uri: image }} style={styles.image} />}
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
  image: {
    width: 200,
    height: 200,
    marginTop: 20,
    borderRadius: 10,
  },
});
