import { StatusBar } from "expo-status-bar";
import { Button, StyleSheet, Text, View } from "react-native";
import { useImages } from "./hooks/useImages";
import ImagePreview from "./components/ImagePreview";

export default function App() {
  const { images, pickImages } = useImages();

  return (
    <View style={styles.container}>
      <Text>Open up App.tsx to start working on your app!</Text>
      <Text>Hello World</Text>
      <StatusBar style="auto" />
      <Button title="Pick an image from your galery" onPress={pickImages} />
      {images && <ImagePreview uri={images[0]} />}
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
