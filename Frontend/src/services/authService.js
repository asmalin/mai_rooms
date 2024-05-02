export async function refreshTokens() {
  try {
    const server_domain = process.env.REACT_APP_SERVER_BASE_URL;
    const response = await fetch(server_domain + "/auth/refresh", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    return response;
  } catch (error) {
    return error;
  }
}

export async function logout() {
  try {
    const server_domain = process.env.REACT_APP_SERVER_BASE_URL;
    const response = await fetch(server_domain + "/logout", {
      method: "GET",
      credentials: "include",
    });

    if (!response.ok) {
      return null;
    }

    const responseText = await response.json();
    return responseText;
  } catch (error) {
    return null;
  }
}

export async function checkAuth() {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const server_domain = process.env.REACT_APP_SERVER_BASE_URL;
    const response = await fetch(server_domain + "/auth/check", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JSON.parse(token)}`,
      },
    });

    return response;
  } catch (error) {
    return error;
  }
}
