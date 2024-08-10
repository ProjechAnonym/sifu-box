export function SocketUrl(
  url: string,
  port: number,
  service: string,
  params: Array<{ [key: string]: string }>
) {
  const newUrl = new URL(url);
  newUrl.protocol = newUrl.protocol === "https:" ? "wss" : "ws";
  newUrl.port = port.toString();
  newUrl.pathname = service;
  params.forEach((param) =>
    Object.entries(param).forEach((item) => newUrl.searchParams.append(...item))
  );
  return newUrl.toString();
}
