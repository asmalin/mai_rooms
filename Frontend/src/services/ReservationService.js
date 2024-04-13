export default async function reserveRoom(lessonForReservation) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch("http://localhost:8080/api/reserve", {
      method: "POST",
      body: JSON.stringify(lessonForReservation),
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JSON.parse(token)}`,
      },
    });

    if (!response.ok) {
      return response;
    }

    const data = await response.json();
    return data;
  } catch (error) {
    return error;
  }
}
