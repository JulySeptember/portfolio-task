function handler(event) {
  var request = event.request;
  var uri = request.uri;

  if (
    !uri.startsWith("/_next/") &&
    !uri.startsWith("/favicon") &&
    !uri.includes(".")
  ) {
    request.uri = uri.replace(/\/?$/, "/index.html");
  }

  return request;
}
