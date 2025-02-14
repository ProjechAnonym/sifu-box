export function sizeCalc(size: number) {
  return size < 1024
    ? `${size} B`
    : size < 1024 * 1024
      ? `${(size / 1024).toFixed(2)} KB`
      : size < 1024 * 1024 * 1024
        ? `${(size / (1024 * 1024)).toFixed(2)} MB`
        : `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}
export function timeCalc(size: number) {
  return size < 60
    ? `${size}分钟`
    : size < 1440
      ? `${(size / 60).toFixed(2)} 小时`
      : size < 60 * 24 * 7
        ? `${(size / (60 * 24)).toFixed(2)}天`
        : `${(size / (60 * 24 * 7)).toFixed(2)}周`;
}
