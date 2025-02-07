export function SocketUrl(
  url: string,
  service: string,
  params: Array<{ [key: string]: string }>
) {
  const newUrl = new URL(url);
  newUrl.protocol = newUrl.protocol === "https:" ? "wss" : "ws";
  newUrl.pathname = service;
  params.forEach((param) =>
    Object.entries(param).forEach((item) => newUrl.searchParams.append(...item))
  );
  return newUrl.toString();
}
