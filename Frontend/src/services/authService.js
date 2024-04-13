export default async function checkAuth() {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return null;
  }

  try {
    const response = await fetch("http://localhost:8080/checkAuth", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JSON.parse(token)}`,
      },
    });

    if (!response.ok) {
      return null;
    }

    const userData = await response.json();
    return userData;
  } catch (error) {
    return null;
  }
}
