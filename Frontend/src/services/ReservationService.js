export async function reserveRoom(lessonForReservation) {
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
      return "error";
    }

    const data = await response.json();
    return data;
  } catch (error) {
    return error;
  }
}

export async function cancelReserve(lessonForCancelReservation) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch(
      "http://localhost:8080/api/cancelReservation",
      {
        method: "POST",
        body: JSON.stringify(lessonForCancelReservation),
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${JSON.parse(token)}`,
        },
      }
    );

    if (!response.ok) {
      return "error";
    }

    const result = await response.json();
    return result;
  } catch (error) {
    return error;
  }
}
