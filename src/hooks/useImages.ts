import { useState } from "react";
import * as ImagePicker from "expo-image-picker";

export function useImages() {
  const [images, setImages] = useState<string[] | null>(null);

  async function pickImages() {
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ["images"],
      allowsEditing: false,
      aspect: [4, 3],
      quality: 1,
    });
    console.log("Image result:", result);

    if (result.canceled) {
      setImages(null);
      return;
    }
    const newImages = result.assets.map((asset) => asset.uri);
    setImages(newImages);
  }

  return { images, pickImages };
}
