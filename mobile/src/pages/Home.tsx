import { Button, StyleSheet, Text, View } from "react-native";
import { useImages } from "../hooks/useImages";
import ImagePreview from "../components/ImagePreview";
import { useNavigation } from "@react-navigation/native";

export function Home() {
  const { images, pickImages } = useImages();
  const navigation = useNavigation();

  return (
    <View style={styles.container}>
      <Text>Welcome — Home screen</Text>
      <Button
        title="Go to Chat"
        onPress={() => navigation.navigate("Chat" as never)}
      />
      <Button title="Pick an image from your gallery" onPress={pickImages} />
      {images && images.length > 0 && <ImagePreview uri={images[0]} />}
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
