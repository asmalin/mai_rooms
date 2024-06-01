export async function refreshTokens() {
  try {
    const response = await fetch("/api/auth/refresh", {
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
    const response = await fetch("/api/auth/logout", {
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
    const response = await fetch("/api/auth/check", {
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
