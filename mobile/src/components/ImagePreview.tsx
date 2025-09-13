import { Image, View } from "react-native";

export default function ImagePreview({ uri }: { uri: string }) {
  return (
    <View style={{ marginTop: 20 }}>
      <Image
        source={{ uri }}
        style={{ width: 200, height: 200, borderRadius: 10 }}
      />
    </View>
  );
}
